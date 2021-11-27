from pymongo import MongoClient
# pprint library is used to make the output look more pretty
from pprint import pprint
import json

username = 'root'
password = 'example'
client = MongoClient(f'mongodb://{username}:{password}@127.0.0.1')

db = client.mydb
persons_col = db.persons
invoices_col = db.invoices
hierarchy_col = db.hierarchy

# Insert objects
with open(r'.\generators\output_data\persons.json', 'r') as f:
    persons = json.load(f)
    persons_col.insert_many(persons)

with open(r'.\generators\output_data\invoices.json', 'r') as f:
    invoices = json.load(f)
    invoices_col.insert_many(invoices)

with open(r'.\generators\output_data\hierarchy.json', 'r') as f:
    hierarchy = json.load(f)
    hierarchy_col.insert_many(hierarchy)

# id = persons.insert_one(person).inserted_id

# # Issue the serverStatus command and print the results
# pprint(id)
