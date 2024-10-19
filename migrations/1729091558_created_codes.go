package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "3owpt1udszv54xg",
			"created": "2024-10-16 15:12:38.985Z",
			"updated": "2024-10-16 15:12:38.985Z",
			"name": "codes",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "ncwrma0b",
					"name": "email",
					"type": "email",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"exceptDomains": [],
						"onlyDomains": []
					}
				},
				{
					"system": false,
					"id": "ae8o7at2",
					"name": "code",
					"type": "text",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				}
			],
			"indexes": [
				"CREATE UNIQUE INDEX ` + "`" + `idx_UZUklhF` + "`" + ` ON ` + "`" + `codes` + "`" + ` (` + "`" + `email` + "`" + `)"
			],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("3owpt1udszv54xg")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
