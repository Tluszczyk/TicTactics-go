// Define collections
collections=[
    "users", 
    "passwordHashes", 
    "userPasswordHashMapping",
    "sessions",
    "userSessionMapping",
    "games",
    "userGameMapping"
];

// Get environment variables
MONGO_INITDB_DATABASE = process.env.MONGO_INITDB_DATABASE;

// Create DB and collection
db = new Mongo().getDB(MONGO_INITDB_DATABASE);
collections.forEach(collection => db.createCollection(collection, { capped: false }));
