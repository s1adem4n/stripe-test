package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("3owpt1udszv54xg")
		if err != nil {
			return err
		}

		// add
		new_expires := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "khufeg0c",
			"name": "expires",
			"type": "date",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), new_expires); err != nil {
			return err
		}
		collection.Schema.AddField(new_expires)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("3owpt1udszv54xg")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("khufeg0c")

		return dao.SaveCollection(collection)
	})
}
