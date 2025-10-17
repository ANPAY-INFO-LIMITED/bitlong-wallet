package rpcclient

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/lightninglabs/lightning-terminal/litrpc"
	"github.com/lightninglabs/taproot-assets/taprpc"
	"github.com/wallet/base"
	"github.com/wallet/service/apiConnect"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func FixAsset(outpoint string, isSpend bool) (string, error) {

	client, clearUp, err := getTaprootAssetsClient()
	if err != nil {
		return "", fmt.Errorf("GetTaprootAssetsClient failed")
	}
	defer clearUp()

	resiveReq := taprpc.AddrReceivesRequest{}
	receives, err := client.AddrReceives(context.Background(), &resiveReq)
	if err != nil {
		return "", err
	}
	for _, receive := range receives.Events {
		if receive.Outpoint == outpoint {
			if receive.HasProof {
				goto list
			}
			return "", fmt.Errorf("not has proof")
		}
	}
	return "", fmt.Errorf("did not hold the asset")
list:
	listReq := taprpc.ListAssetRequest{
		IncludeSpent: true,
	}
	assets, err := client.ListAssets(context.Background(), &listReq)
	if err != nil {
		return "", fmt.Errorf("GetListAssets failed")
	}
	for _, asset := range assets.Assets {
		if asset.ChainAnchor.AnchorOutpoint == outpoint {
			goto tx
		}
	}
	return "", fmt.Errorf("did not find the asset in spent list")

tx:
	coon, clearUp, err := apiConnect.GetConnection("litd", false)
	if err != nil {
		return "", err
	}
	defer clearUp()
	litdClient := litrpc.NewProxyClient(coon)
	stopReq := &litrpc.StopDaemonRequest{}

	_, err = litdClient.StopDaemon(context.Background(), stopReq)
	if err != nil {
		return "", errors.New("litd is not stop")
	}
	for i := 0; i < 10; i++ {
		getinfoReq := taprpc.GetInfoRequest{}
		_, err := client.GetInfo(context.Background(), &getinfoReq)
		if err != nil {
			fmt.Println("tapd is stop")
			time.Sleep(time.Second * 5)
			goto dbSet
		}
	}
	return "", errors.New("tapd is not stop")

dbSet:
	dbPath := filepath.Join(base.Configure("tapd"), "data", base.NetWork, "tapd.db")
	fmt.Printf("FixAsset start, dbPath: %s\n, outpoint: %s\n", dbPath, outpoint)
	db, err := sql.Open("sqlite", dbPath+"?_busy_timeout=500000")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	op := strings.Split(outpoint, ":")
	b, err := hex.DecodeString(op[0])
	if err != nil {
		return "", err
	}
	for i := 0; i < len(b)/2; i++ {
		temp := b[i]
		b[i] = b[len(b)-i-1]
		b[len(b)-i-1] = temp
	}
	outBytes := make([]byte, 36)
	copy(outBytes[0:32], b)
	port, err := strconv.Atoi(op[1])
	outBytes[32] = 0x00 + byte(port)
	out := hex.EncodeToString(outBytes)
	var Spend int
	if isSpend {
		Spend = 1
	}
	lows, err := db.Exec(test, out, Spend)
	if err != nil {
		return "", err
	}
	affectRows, _ := lows.RowsAffected()
	fmt.Printf("FixAsset start end, affected rows: %d\n", affectRows)
	return "", nil
}

const test = `
	update  assets 
	set spent = $2
	WHERE anchor_utxo_id = (
		SELECT utxo_id 
		FROM managed_utxos 
		WHERE lower(hex(outpoint)) = $1
	);
	`

func CheckTapdDb() (bool, error) {
	dbPath := filepath.Join(base.Configure("tapd"), "data", base.NetWork, "tapd.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return true, nil
	}
	fmt.Printf("CheckTapdDb start, dbPath: %s\n", dbPath)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	defer db.Close()

	_, err = db.Exec(checkpoint)
	if err != nil {
		return false, err
	}

	var name string
	err = db.QueryRow(checktabl, "multiverse_roots").Scan(&name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("table not exist")
		_, err = db.Exec(createrootstable)
		if err != nil {
			return false, err
		}
	}
	fmt.Printf("CheckTapdDb start is false end\n")
	return true, nil
}

const checkpoint = `
	PRAGMA wal_checkpoint(FULL);
`
const checktabl = `
SELECT name 
FROM sqlite_master 
WHERE type='table' AND name= ?
`
const createrootstable = `
create table multiverse_roots
(
    id             INTEGER
        primary key,
    namespace_root VARCHAR not null
        unique
        references mssmt_roots
            deferrable initially deferred,
    proof_type     TEXT    not null,
    check (proof_type IN ('issuance', 'transfer'))
);
INSERT INTO multiverse_roots ('namespace_root', 'proof_type') VALUES ('multiverse-issuance', 'issuance');
INSERT INTO multiverse_roots ('namespace_root', 'proof_type') VALUES ('multiverse-transfer', 'transfer');
create table multiverse_leaves
(
    id                  INTEGER
        primary key,
    multiverse_root_id  BIGINT  not null
        references multiverse_roots,
    asset_id            BLOB,
    group_key           BLOB,
    leaf_node_key       BLOB    not null,
    leaf_node_namespace VARCHAR not null,
    check ((asset_id IS NOT NULL AND group_key IS NULL) OR
           (asset_id IS NULL AND group_key IS NOT NULL)),
    check (LENGTH(group_key) = 32),
    check (length(asset_id) = 32)
);

create unique index multiverse_leaves_unique
    on multiverse_leaves (leaf_node_key, leaf_node_namespace);

`
