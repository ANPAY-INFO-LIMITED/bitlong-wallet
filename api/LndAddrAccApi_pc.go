package api

import (
	"context"
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wallet/base"
	"github.com/wallet/service/apiConnect"
	"github.com/wallet/service/rpcclient"
)

func (s *AddrStore) AllAddresses2(bucket string) ([]*Addr, error) {
	var Addrs []*Addr
	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.ForEach(func(k, v []byte) error {
			var u Addr
			err := json.Unmarshal(v, &u)
			if err != nil {
				return err
			}
			Addrs = append(Addrs, &u)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return Addrs, nil
}

func PcGetNewAddress_P2TR() (*Addr, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()

	lc := lnrpc.NewLightningClient(conn)
	request := &lnrpc.NewAddressRequest{
		Type: lnrpc.AddressType_TAPROOT_PUBKEY,
	}
	response, err := lc.NewAddress(context.Background(), request)
	if err != nil {
		return nil, errors.Wrap(err, "lc.NewAddress")
	}
	return &Addr{
		Name:           "default",
		Address:        response.Address,
		Balance:        0,
		AddressType:    lnrpc.AddressType_TAPROOT_PUBKEY.String(),
		DerivationPath: AddressTypeToDerivationPath(lnrpc.AddressType_TAPROOT_PUBKEY.String()),
		IsInternal:     false,
	}, nil
}

func PcGetNewAddress_P2WKH() (*Addr, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.NewAddressRequest{
		Type: lnrpc.AddressType_WITNESS_PUBKEY_HASH,
	}
	response, err := client.NewAddress(context.Background(), request)
	if err != nil {
		return nil, errors.Wrap(err, "lc.NewAddress")
	}
	return &Addr{
		Name:           "default",
		Address:        response.Address,
		Balance:        0,
		AddressType:    lnrpc.AddressType_WITNESS_PUBKEY_HASH.String(),
		DerivationPath: AddressTypeToDerivationPath(lnrpc.AddressType_WITNESS_PUBKEY_HASH.String()),
		IsInternal:     false,
	}, nil
}

func PcGetNewAddress_NP2WKH() (*Addr, error) {
	conn, clearUp, err := apiConnect.GetConnection("lnd", false)
	if err != nil {
		return nil, errors.Wrap(err, "apiConnect.GetConnection")
	}
	defer clearUp()
	client := lnrpc.NewLightningClient(conn)
	request := &lnrpc.NewAddressRequest{
		Type: lnrpc.AddressType_NESTED_PUBKEY_HASH,
	}
	response, err := client.NewAddress(context.Background(), request)
	if err != nil {
		return nil, errors.Wrap(err, "lc.NewAddress")
	}
	return &Addr{
		Name:           "default",
		Address:        response.Address,
		Balance:        0,
		AddressType:    lnrpc.AddressType_NESTED_PUBKEY_HASH.String(),
		DerivationPath: AddressTypeToDerivationPath(lnrpc.AddressType_NESTED_PUBKEY_HASH.String()),
		IsInternal:     false,
	}, nil
}

func PcGetNewAddress_P2TR_Example() (*Addr, error) {
	address := "bc1pxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	return &Addr{
		Name:           "default",
		Address:        address,
		Balance:        0,
		AddressType:    lnrpc.AddressType_TAPROOT_PUBKEY.String(),
		DerivationPath: AddressTypeToDerivationPath(lnrpc.AddressType_TAPROOT_PUBKEY.String()),
		IsInternal:     false,
	}, nil
}

func PcGetNewAddress_P2WKH_Example() (*Addr, error) {
	address := "bc1qxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	return &Addr{
		Name:           "default",
		Address:        address,
		Balance:        0,
		AddressType:    lnrpc.AddressType_WITNESS_PUBKEY_HASH.String(),
		DerivationPath: AddressTypeToDerivationPath(lnrpc.AddressType_WITNESS_PUBKEY_HASH.String()),
		IsInternal:     false,
	}, nil
}

func PcGetNewAddress_NP2WKH_Example() (*Addr, error) {
	address := "3xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	return &Addr{
		Name:           "default",
		Address:        address,
		Balance:        0,
		AddressType:    lnrpc.AddressType_NESTED_PUBKEY_HASH.String(),
		DerivationPath: AddressTypeToDerivationPath(lnrpc.AddressType_NESTED_PUBKEY_HASH.String()),
		IsInternal:     false,
	}, nil
}

func PcStoreAddr(name string, address string, balance int, addressType string, derivationPath string, isInternal bool) error {
	err := InitAddrDB()
	if err != nil {
		return errors.Wrap(err, "InitAddrDB")
	}
	path := filepath.Join(base.QueryConfigByKey("dirpath"), "phone.db")
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return errors.Wrap(err, "bolt.Open")
	}
	defer func(db *bolt.DB) {
		err := db.Close()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "db.Close"))
		}
	}(db)
	s := &AddrStore{DB: db}
	err = s.CreateOrUpdateAddr("addresses", &Addr{
		Name:           name,
		Address:        address,
		Balance:        balance,
		AddressType:    addressType,
		DerivationPath: derivationPath,
		IsInternal:     isInternal,
	})
	if err != nil {
		return errors.Wrap(err, "s.CreateOrUpdateAddr")
	}
	return nil
}

func PcRemoveAddr(address string) error {
	err := InitAddrDB()
	if err != nil {
		return errors.Wrap(err, "InitAddrDB")
	}
	path := filepath.Join(base.QueryConfigByKey("dirpath"), "phone.db")
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return errors.Wrap(err, "bolt.Open")
	}
	defer func(db *bolt.DB) {
		err := db.Close()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "db.Close"))
		}
	}(db)
	s := &AddrStore{DB: db}
	_, err = s.ReadAddr("addresses", address)
	if err != nil {
		return errors.Wrap(err, "s.ReadAddr")
	}
	err = s.DeleteAddr("addresses", address)
	if err != nil {
		return errors.Wrap(err, "s.DeleteAddr")
	}
	return nil
}

func PcQueryDbAddr(address string) (*Addr, error) {
	err := InitAddrDB()
	if err != nil {
		return nil, errors.Wrap(err, "InitAddrDB")
	}
	path := filepath.Join(base.QueryConfigByKey("dirpath"), "phone.db")
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, errors.Wrap(err, "bolt.Open")
	}
	defer func(db *bolt.DB) {
		err := db.Close()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "db.Close"))
		}
	}(db)
	s := &AddrStore{DB: db}
	addr, err := s.ReadAddr("addresses", address)
	if err != nil {
		return nil, errors.Wrap(err, "s.ReadAddr")
	}
	return addr, nil
}

func PcQueryAllAddr() ([]*Addr, error) {
	err := InitAddrDB()
	if err != nil {
		return nil, errors.Wrap(err, "InitAddrDB")
	}
	path := filepath.Join(base.QueryConfigByKey("dirpath"), "phone.db")
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, errors.Wrap(err, "bolt.Open")
	}
	defer func(db *bolt.DB) {
		err := db.Close()
		if err != nil {
			logrus.Errorln(errors.Wrap(err, "db.Close"))
		}
	}(db)
	s := &AddrStore{DB: db}
	addresses, err := s.AllAddresses2("addresses")
	if err != nil {
		return nil, errors.Wrap(err, "s.AllAddresses2")
	}
	return addresses, nil
}

func PcGetNonZeroBalanceAddresses() ([]*Addr, error) {
	listAddrResp, err := rpcclient.ListAddresses()
	if err != nil {
		return nil, errors.Wrap(err, "rpcclient.ListAddresses")
	}
	var addrs []*Addr
	listAddrs := listAddrResp.GetAccountWithAddresses()
	for _, accWithAddr := range listAddrs {
		addresses := accWithAddr.Addresses
		for _, address := range addresses {
			if address.Balance != 0 {
				addrs = append(addrs, &Addr{
					Name:           accWithAddr.Name,
					Address:        address.Address,
					Balance:        int(address.Balance),
					AddressType:    accWithAddr.AddressType.String(),
					DerivationPath: accWithAddr.DerivationPath,
					IsInternal:     address.IsInternal,
				})
			}
		}
	}
	return addrs, nil
}

func PcUpdateAllAddressesByGNZBA() error {
	listAddrResp, err := rpcclient.ListAddresses()
	if err != nil {
		return errors.Wrap(err, "rpcclient.ListAddresses")
	}
	var addresses []string
	listAddrs := listAddrResp.GetAccountWithAddresses()
	allAddr, err := QueryAllAddrAndGetResponse()
	if err != nil {
		return errors.Wrap(err, "QueryAllAddrAndGetResponse")
	}
	// @dev: Update allAddr balance
	err = UpdateAllAddrByAccountWithAddresses(allAddr, &listAddrs)
	if err != nil {
		return errors.Wrap(err, "UpdateAllAddrByAccountWithAddresses")
	}
	// @dev: UpdateNoneZeroAddress
	for _, accWithAddr := range listAddrs {
		if accWithAddr.Name != "default" {
			continue
		}
		_addresses := accWithAddr.Addresses
		for _, _address := range _addresses {
			// @dev: remove is_internal check
			if _address.Balance != 0 {
				// @dev: Store
				err = PcStoreAddr(accWithAddr.Name, _address.Address, int(_address.Balance), accWithAddr.AddressType.String(), accWithAddr.DerivationPath, _address.IsInternal)
				if err != nil {
					return errors.Wrap(err, "PcStoreAddr")
				}
				addresses = append(addresses, _address.Address)
			}
		}
	}
	return nil
}
