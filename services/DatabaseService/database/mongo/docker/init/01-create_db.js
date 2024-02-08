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
SESSION_TTL_SECONDS = parseInt(process.env.SESSION_TTL_SECONDS, 10);

// Create DB and collection
db = new Mongo().getDB(MONGO_INITDB_DATABASE);
collections.forEach(collection => db.createCollection(collection, { capped: false }));

db.sessions.createIndex({ "createdAt": 1 }, { expireAfterSeconds: SESSION_TTL_SECONDS });
db.userSessionMapping.createIndex({ "createdAt": 1 }, { expireAfterSeconds: SESSION_TTL_SECONDS });