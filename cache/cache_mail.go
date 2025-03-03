package cache

import (
	"context"
	"fmt"
	"r3/types"
	"sync"

	"github.com/jackc/pgx/v5"
)

var (
	mail_mx          sync.RWMutex
	mailAccountIdMap map[int32]types.MailAccount
)

func GetMailAccountMap() map[int32]types.MailAccount {
	mail_mx.RLock()
	defer mail_mx.RUnlock()

	return mailAccountIdMap
}

func GetMailAccount(id int32, mode string) (types.MailAccount, error) {
	mail_mx.RLock()
	defer mail_mx.RUnlock()

	ma, exists := mailAccountIdMap[id]
	if !exists || mode != ma.Mode {
		return ma, fmt.Errorf("mail account with ID %d does not exist for mode '%s'", id, mode)
	}
	return ma, nil
}

func GetMailAccountAny(mode string) (types.MailAccount, error) {
	mail_mx.RLock()
	defer mail_mx.RUnlock()

	for _, ma := range mailAccountIdMap {
		if mode == ma.Mode {
			return ma, nil
		}
	}
	return types.MailAccount{}, fmt.Errorf("no mail account is available for mode '%s'", mode)
}

func GetMailAccountsExist() bool {
	mail_mx.RLock()
	defer mail_mx.RUnlock()

	return len(mailAccountIdMap) != 0
}

func LoadMailAccountMap_tx(ctx context.Context, tx pgx.Tx) error {

	rows, err := tx.Query(ctx, `
		SELECT id, oauth_client_id, name, mode, auth_method, username,
			password, start_tls, send_as, host_name, host_port, comment
		FROM instance.mail_account
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	mail_mx.Lock()
	defer mail_mx.Unlock()

	mailAccountIdMap = make(map[int32]types.MailAccount)
	for rows.Next() {
		var ma types.MailAccount

		if err := rows.Scan(&ma.Id, &ma.OauthClientId, &ma.Name, &ma.Mode,
			&ma.AuthMethod, &ma.Username, &ma.Password, &ma.StartTls,
			&ma.SendAs, &ma.HostName, &ma.HostPort, &ma.Comment); err != nil {

			return err
		}
		mailAccountIdMap[ma.Id] = ma
	}
	return nil
}
