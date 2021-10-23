from typing import Optional
from fastapi import FastAPI
import couchdb
from snorkel import Snorkel
import os


DB_NAME = "packer"

snorkel = Snorkel("packer", "__snorkel-relay__")
snorkel.add_str_field("event")
snorkel.add_str_field("path")
snorkel.add_str_field("error_name")
snorkel.add_str_field("error")
snorkel.add_int_field("elapsed")
snorkel.add_str_field("extras")


def connect_db(*, url: str):
    server = couchdb.Server(url)
    return server[DB_NAME] if DB_NAME in server else server.create(DB_NAME)


db = connect_db(url=os.environ.get('DB_URL'))
app = FastAPI()


@app.get("/")
def read_root():
    return {"Hello": "World"}


@app.get("/items/{item_id}")
def read_item(item_id: str, q: Optional[str] = None):
    if item_id in db:
        return db[item_id]
    db[item_id] = {"id": item_id, "q": q}
    snorkel.write({
        "event": "item-created",
        "path": f"/items/{item_id}",
        "extras": json.dumps({"item_id": item_id})
    })

    return db[item_id]


@app.get("/items")
def read_item_list():
    return [item for item in db]
