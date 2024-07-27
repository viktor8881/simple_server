package user

type Model struct {
	ID    string `bson:"_id,omitempty"`
	Name  string `bson:"name"`
	Email string `bson:"value"`
}
