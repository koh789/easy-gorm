# easy-gorm


[![BuildStatus](https://github.com/koh789/easy-gorm/actions/workflows/test.yml/badge.svg)](https://github.com/koh789/easy-gorm/actions/workflows/test.yml)


## Quick Start


Using the following user table as an example, embed CRUDClient within the target Client implementation class (userClientImpl).
And embed the CRUD interface in the UserClient interface.

```go

type UserClient interface {
	egorm.CRUD[User, UserPK, Users]
}

type userClientImpl struct {
	*egorm.CRUDClient[User, UserPK, Users]
}

func NewUserClientImpl(db *gorm.Db)UserClient{
	return &userClientImpl{CRUDClient: db}
}

type User struct {
	PK   UserPK
	Name string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserPK struct {
	Email string `gorm:"primaryKey;autoIncrement:false"`
}

type Users []User

```

Then, as shown below, a CRUD call can be made to manipulate the target table.


```go

func example(){
	var db *gorm.DB = // create gorm.DB instance 
	
	client := NewUserClientImpl(db)
	pk := UserPK{Email: "dummy@email"}
	// find by pk
	client.FindByID(pk)
	// find by pk's slice
	client.FindByIDs([]UserPK{pk})
	// find all
	client.FindAll()
	user := User{PK: pk, Name: "dummy"}
	// insert on duplicate key update
	client.Save(&user)
	// insert on duplicate key update
	client.SaveAll(Users{user})
}

```