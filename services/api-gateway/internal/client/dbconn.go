package client

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// Database credentials come only from the environment (docker-compose / .env); nothing is hardcoded here.
//
//	SUPER_ADMIN_DB_USER / SUPER_ADMIN_DB_PASSWORD  shared root account on postgres_apps (tried first when set)
//	AUTH_DB_* / USER_DB_* / ATTENDANCE_DB_*        per-service USER, PASSWORD and NAME

type dbCreds struct {
	user string
	pass string
	name string
}

// lastGoodConn remembers the connection string that last worked per database, so requests do not
// walk through unreachable hosts (DNS failures cost seconds each) every time.
var lastGoodConn sync.Map

func envOr(keys []string, def string) string {
	for _, k := range keys {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			return v
		}
	}
	return def
}

func serviceCreds(prefix, defUser, defName string) dbCreds {
	return dbCreds{
		user: envOr([]string{prefix + "_DB_USER", "DB_USER"}, defUser),
		pass: envOr([]string{prefix + "_DB_PASSWORD", "DB_PASSWORD"}, ""),
		name: envOr([]string{prefix + "_DB_NAME", "DB_NAME"}, defName),
	}
}

func dbHostCandidates() []string {
	hosts := []string{os.Getenv("DB_HOST"), "postgres_apps", "localhost", "127.0.0.1", "host.docker.internal"}
	seen := map[string]bool{}
	out := []string{}
	for _, h := range hosts {
		h = strings.TrimSpace(h)
		if h != "" && !seen[h] {
			seen[h] = true
			out = append(out, h)
		}
	}
	return out
}

// quoteConnValue escapes a value for a libpq key='value' connection string.
func quoteConnValue(v string) string {
	return "'" + strings.ReplaceAll(strings.ReplaceAll(v, `\`, `\\`), `'`, `\'`) + "'"
}

func connString(host string, c dbCreds) string {
	return fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable connect_timeout=3",
		host, quoteConnValue(c.user), quoteConnValue(c.pass), quoteConnValue(c.name))
}

// buildConnStrings lists every candidate connection for the service database: the super admin account
// first (when configured), then the service account, on each host. A cached working one comes first.
func buildConnStrings(svc dbCreds) []string {
	superUser := envOr([]string{"SUPER_ADMIN_DB_USER"}, "super_admin_apps")
	superPass := envOr([]string{"SUPER_ADMIN_DB_PASSWORD"}, "")

	var conns []string
	seen := map[string]bool{}
	add := func(c string) {
		if !seen[c] {
			seen[c] = true
			conns = append(conns, c)
		}
	}
	if cached, ok := lastGoodConn.Load(svc.name); ok {
		add(cached.(string))
	}
	for _, h := range dbHostCandidates() {
		if superPass != "" {
			add(connString(h, dbCreds{user: superUser, pass: superPass, name: svc.name}))
		}
		if svc.pass != "" {
			add(connString(h, svc))
		}
	}
	return conns
}

func rememberGoodConn(dbName, conn string) {
	lastGoodConn.Store(dbName, conn)
}

// openFirstReachable returns the first candidate that answers a ping.
func openFirstReachable(dbName string, conns []string) (*sql.DB, error) {
	var lastErr error = fmt.Errorf("tidak ada kredensial database %s yang dikonfigurasi (cek variabel environment *_DB_PASSWORD)", dbName)
	for _, conn := range conns {
		db, err := sql.Open("postgres", conn)
		if err != nil {
			lastErr = err
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err = db.PingContext(ctx)
		cancel()
		if err == nil {
			rememberGoodConn(dbName, conn)
			return db, nil
		}
		lastErr = err
		db.Close()
	}
	log.Printf("[DB Error] Tidak dapat terhubung ke %s: %v", dbName, lastErr)
	return nil, lastErr
}

func authDBCreds() dbCreds {
	return serviceCreds("AUTH", "user_lopiq_auth", "db_lopiq_auth")
}

func userDBCreds() dbCreds {
	return serviceCreds("USER", "user_lopiq_user", "db_lopiq_user")
}

func attendanceDBCreds() dbCreds {
	return serviceCreds("ATTENDANCE", "user_lopiq_attendance", "db_lopiq_attendance")
}

func openAuthDB() (*sql.DB, error) {
	c := authDBCreds()
	return openFirstReachable(c.name, buildConnStrings(c))
}
