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
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func FixAsset(outpoint string) (string, error) {

	client, clearUp, err := getTaprootAssetsClient()
	if err != nil {
		return "", fmt.Errorf("GetTaprootAssetsClient failed")
	}
	defer clearUp()

	//接收列表中是否存在该资产
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
	//已使用列表中是否存在该资产
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

	//交易列表中是否存在该资产
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
	//如果存在，则将spent置为0
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
	lows, err := db.Exec(test, out)
	if err != nil {
		return "", err
	}
	affectRows, _ := lows.RowsAffected()
	fmt.Printf("FixAsset start end, affected rows: %d\n", affectRows)
	return "", nil
}

const test = `
	update  assets 
	set spent = 0
	WHERE anchor_utxo_id = (
		SELECT utxo_id 
		FROM managed_utxos 
		WHERE lower(hex(outpoint)) = $1
	);
	`
