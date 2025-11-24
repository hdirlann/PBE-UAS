package postgre

type Permission struct {
    ID          string `db:"id" json:"id"`
    Name        string `db:"name" json:"name"`
    Resource    string `db:"resource" json:"resource"`
    Action      string `db:"action" json:"action"`
    Description string `db:"description" json:"description"`
}
