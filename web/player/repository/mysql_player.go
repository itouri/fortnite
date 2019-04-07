package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/itouri/fortnite/web/domain"
	"github.com/itouri/fortnite/web/player"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00"
)

type mysqlPlayerRepository struct {
	Conn *sql.DB
}

func (m *mysqlPlayerRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*domain.Player, error) {

	rows, err := m.Conn.QueryContext(ctx, query, args...)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()
	result := make([]*domain.Player, 0)
	for rows.Next() {
		t := new(domain.Player)
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.IconPath,
			&t.CoverImgPath,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func NewMysqlPlayerRepository(Conn *sql.DB) player.Repository {

	return &mysqlPlayerRepository{Conn}
}

func (m *mysqlPlayerRepository) Fetch(ctx context.Context, cursor string, num int64) ([]*domain.Player, string, error) {

	query := `SELECT id, name, icon_path, cover_img_path, created_at, updated,at
				FROM players WHERE created_at > ? ORDER BY created_at LIMIT ? `
	decodedCursor, err := DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}
	res, err := m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}
	nextCursor := ""
	//LEARN
	if len(res) == int(num) {
		nextCursor = EncodeCursor(res[len(res)-1].CreatedAt)
	}
	return res, nextCursor, err
}

func (m *mysqlPlayerRepository) GetByID(ctx context.Context, id int64) (*domain.Player, error) {

	//FIXME アスタリスクではダメなの？
	query := `SELECT id, name, icon_path, cover_img_path, created_at, updated_at
				FROM players WHERE ID = ?`
	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	p := &domain.Player{}
	if len(list) > 0 {
		p = list[0]
	} else {
		//FIXME エラーにするほどなのかな？
		return nil, domain.ErrNotFound
	}

	return p, nil
}

func (m *mysqlPlayerRepository) Store(ctx context.Context, p *domain.Player) error {

	query := `INSERT players SET name=?, icon_path=?, cover_img_path=? created_at=? updated_at=?`
	//LEARN
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	logrus.Debug("Created At: ", p.CreatedAt)
	//LEARN
	res, err := stmt.ExecContext(ctx, p.Name, p.IconPath, p.CoverImgPath, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = lastID
	return nil
}

func (m *mysqlPlayerRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM players WHERE id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAfected != 1 {
		//FIXME domain/errors.go に定義しないの?
		err = fmt.Errorf("Weird Behaviour. Toal Affected %d", rowsAfected)
		return err
	}

	return nil
}

func (m *mysqlPlayerRepository) Update(ctx context.Context, p *domain.Player) error {
	//CONSIDER なんで set が小文字になった?
	query := `UPDATE players set name=?, icon_path=?, cover_img_path=?, update_at=? WHERE ID = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, p.Name, p.IconPath, p.CoverImgPath, p.UpdatedAt)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect != 1 {
		err = fmt.Errorf("Weird Behaviour. Total Affected: %d", affect)
		return err
	}
	return nil
}

//LEARN 何をしている？
func DecodeCursor(encodedTime string) (time.Time, error) {
	b, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	timeString := string(b)
	t, err := time.Parse(timeFormat, timeString)
	return t, err
}

//LEARN 何をしている？
func EncodeCursor(t time.Time) string {
	timeString := t.Format(timeFormat)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}
