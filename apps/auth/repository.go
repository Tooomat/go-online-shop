package auth

import (
	"context"
	"database/sql"

	"github.com/Tooomat/go-online-shop/infrastructure/response"
	"github.com/jmoiron/sqlx"
)

type Repository interface { //contract
	CreatedAuth(ctx context.Context, model AuthEntity) (err error)
	GetAuthByEmail(ctx context.Context, email string) (model AuthEntity, err error)
	CekSuperAdmin(ctx context.Context) (count int, err error)
}

type repository struct {
	db *sqlx.DB
}

func newRepository(db *sqlx.DB) Repository {
	return repository{
		db: db,
	}
}

func (r repository) CekSuperAdmin(ctx context.Context) (count int, err error) {
	query := `
		SELECT COUNT(role)
		FROM auth
		WHERE role = 'super_admin'
	`

	if err = r.db.GetContext(ctx, &count, query); err != nil {
		return
	}

	return
}

// step3(register): memasukkan inputan (mapping) yang telah dimasukkan distruct ke database (lewat query)
func (r repository) CreatedAuth(ctx context.Context, model AuthEntity) (err error) {
	query := `
		INSERT INTO auth(
			email, password, role, created_time, update_time, public_id
		) VALUES (
		 	:email, :password, :role, :created_time, :update_time, :public_id
		)
	`

	//consumed 1 pool
	stmt, err := r.db.PrepareNamedContext(ctx, query) //cocok untuk berulang kali insert banyak user
	// fmt.Println(stmt)
	if err != nil {
		return
	}
	defer stmt.Close() //setelah consumed di close

	//GetContex() -> ambil 1 baris langsung dari struct
	//SelectContex() -> ambil banyak baris → langsung ke slice of struct.
	_, err = stmt.ExecContext(ctx, &model) //Untuk INSERT / UPDATE / DELETE (tidak ada return rows).
	//QueryxContext() -> scan() manual
	//NamedExecContext() -> cocok untuk sekali execute
	return
}

// mencari email apakah ada yang sama atau tidak
func (r repository) GetAuthByEmail(ctx context.Context, email string) (model AuthEntity, err error) {
	query := `
		SELECT id, email, password, role, created_time, update_time, public_id
		FROM auth
		WHERE email=?
	`

	err = r.db.GetContext(ctx, &model, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			err = response.ErrNotFound
			return
		}
		return
	}
	return
}
