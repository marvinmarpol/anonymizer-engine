package mask

type Mask struct {
	tableName struct{} `pg:"mask" json:"-"`
	Token     string   `pg:"token,notnull,pk,type:varchar" json:"token"`
	Hash      string   `pg:"hash,notnull,unique,type:varchar" json:"hash"`
	Key       string   `pg:"key,notnull,type:varchar" json:"key"`
	Cypher    string   `pg:"cypher,notnull,type:varchar" json:"cypher"`
	CreatedAt string   `pg:"created_at,notnull,type:timestamptz,default:now()" json:"-"`
	UpdatedAt string   `pg:"updated_at,notnull,type:timestamptz,default:now()" json:"-"`
	RotatedAt string   `pg:"rotated_at,notnull,type:timestamptz,default:now()" json:"-"`
}
