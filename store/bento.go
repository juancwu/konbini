package store

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"time"
)

// Bento represents a bento object that has a slice with all the entries as well.
type Bento struct {
	Id        string
	Name      string
	OwnerId   string
	PubKey    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Delete removes the bento from the database.
//
// You must call tx.Commit for it to take effect.
func (b *Bento) Delete(tx *sql.Tx) (sql.Result, error) {
	return tx.Exec("DELETE FROM bentos WHERE id = $1;", b.Id)
}

// VerifySignature verifies the given hex encoded signature with the given challenge.
func (b *Bento) VerifySignature(signature string, challenge string) error {
	// decode signature
	decodedSignature, err := hex.DecodeString(signature)
	if err != nil {
		return err
	}
	// hashed challenge
	decodedChallenge, err := hex.DecodeString(challenge)
	if err != nil {
		return err
	}
	hashed := sha256.Sum256(decodedChallenge)
	// decode pem encoded public key
	block, _ := pem.Decode([]byte(b.PubKey))
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	pubKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return errors.New("Invalid public key.")
	}
	// verify signature
	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], decodedSignature)
}

// Rename will rename the bento by updating the database record.
// If the new name is the same as the current name, it won't perform a database query.
// After a successful rename, it will auto assigned the new name to the Bento instance.
func (b *Bento) Rename(newName string) error {
	// no need to update
	if b.Name == newName {
		return nil
	}
	_, err := db.Exec("UPDATE bentos SET name = $1 WHERE id = $2;", newName, b.Id)
	if err != nil {
		return err
	}
	b.Name = newName
	return nil
}

// NewBentoTx will create and save a new bento into the database with the given information.
// This method will return an error if there is another bento with the same name from the same user.
// All bentos belonging to one user should have unique names.
//
// IMPORTANT: You must call tx.Commit() afterwards for changes to take effect.
func NewBentoTx(tx *sql.Tx, name, ownerId, pubKey string) (*Bento, error) {
	row := db.QueryRow(
		"INSERT INTO bentos (name, owner_id, pub_key) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at;",
		name,
		ownerId,
		pubKey,
	)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	bento := Bento{
		Name:    name,
		OwnerId: ownerId,
		PubKey:  pubKey,
	}
	err = row.Scan(
		&bento.Id,
		&bento.CreatedAt,
		&bento.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &bento, nil
}

// GetBentoWithId retrieves a bento with the given id.
func GetBentoWithId(id string) (*Bento, error) {
	row := db.QueryRow(
		"SELECT name, owner_id, pub_key, created_at, updated_at FROM bentos WHERE id = $1;",
		id,
	)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	bento := Bento{
		Id: id,
	}
	err = row.Scan(
		&bento.Name,
		&bento.OwnerId,
		&bento.PubKey,
		&bento.CreatedAt,
		&bento.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &bento, nil
}
