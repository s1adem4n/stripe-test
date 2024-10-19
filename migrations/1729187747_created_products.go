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
			"id": "nqpzhu9bq409ffk",
			"created": "2024-10-17 17:55:47.261Z",
			"updated": "2024-10-17 17:55:47.261Z",
			"name": "products",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "7ve8llmm",
					"name": "name",
					"type": "text",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": null,
						"max": null,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "yiucryjx",
					"name": "price",
					"type": "number",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"min": 1,
						"max": null,
						"noDecimal": true
					}
				},
				{
					"system": false,
					"id": "dpqbuh6i",
					"name": "description",
					"type": "editor",
					"required": false,
					"presentable": false,
					"unique": false,
					"options": {
						"convertUrls": false
					}
				}
			],
			"indexes": [],
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

		collection, err := dao.FindCollectionByNameOrId("nqpzhu9bq409ffk")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
