package repositories

import (
	"context"
	"database/sql"

	"github.com/Bayan2019/hackathon-2025-police-volunteers/repositories/database"
	"github.com/Bayan2019/hackathon-2025-police-volunteers/views"
)

type UsersRepository struct {
	Conn *sql.DB
	DB   *database.Queries
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{
		Conn: db,
		DB:   database.New(db),
	}
}

func (ur *UsersRepository) Create(ctx context.Context, cup database.CreateUserParams) (int64, error) {
	tx, err := ur.Conn.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	qtx := ur.DB.WithTx(tx)

	id, err := qtx.CreateUser(ctx, cup)
	if err != nil {
		return 0, err
	}

	err = qtx.AddRole2User(ctx, database.AddRole2UserParams{
		UserID: id,
		RoleID: 1,
	})
	if err != nil {
		return 0, err
	}

	return id, tx.Commit()
}

func (ur *UsersRepository) UpdateProfile(ctx context.Context, id int64, upr views.UpdateProfileRequest) error {
	tx, err := ur.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := ur.DB.WithTx(tx)

	err = qtx.UpdateUser(ctx, database.UpdateUserParams{
		ID:   id,
		Name: upr.Name,
		// Ph:       upr.Email,
		DateOfBirth: upr.DateOfBirth,
		Phone:       upr.Phone,
		Iin:         upr.Iin,
	})
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (ur *UsersRepository) Update(ctx context.Context, id int64, uur views.UpdateUserRequest) error {
	tx, err := ur.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := ur.DB.WithTx(tx)

	// isPolice := false
	// for _, role_id := range uur.RoleIds {
	// 	if role_id == 2 {
	// 		isPolice = true
	// 		break
	// 	}
	// }

	// if !isPolice {
	// 	uur.CenterId = 0
	// }

	err = qtx.UpdateUser(ctx, database.UpdateUserParams{
		ID:   id,
		Name: uur.Name,
		// Email:       uur.Email,
		DateOfBirth: uur.DateOfBirth,
		Phone:       uur.Phone,
		Iin:         uur.Iin,
		// CenterID:    uur.CenterId,
	})
	if err != nil {
		return err
	}

	err = qtx.RemoveRolesOfUser(ctx, id)
	if err != nil {
		return err
	}

	for _, role_id := range uur.RoleIds {
		err = qtx.AddRole2User(ctx, database.AddRole2UserParams{
			UserID: id,
			RoleID: role_id,
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (ur *UsersRepository) Delete(ctx context.Context, id int64) error {
	tx, err := ur.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := ur.DB.WithTx(tx)

	err = qtx.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
