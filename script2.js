const sqlite3 = require('sqlite3').verbose();
const { v4: uuidv4 } = require('uuid');
const fs = require('fs');
const path = require('path');

console.log('🚀 Starting database population script:', new Date().toISOString());

const dbPath = './backend/database.sqlite';
const schemaPath = './backend/internal/db/schema.sql';

// Ensure the backend directory exists
if (!fs.existsSync('./backend')) {
    fs.mkdirSync('./backend');
    console.log('📁 Created backend directory:', new Date().toISOString());
}

// Connect to the SQLite database (it will be created if it doesn't exist)
const db = new sqlite3.Database(dbPath, (err) => {
    if (err) {
        console.error('❌ Error connecting to the database:', err.message, new Date().toISOString());
        process.exit(1);
    } else {
        console.log('✅ Connected to the SQLite database:', new Date().toISOString());
        initializeDatabase();
    }
});

function initializeDatabase() {
    console.log('🏗️ Checking database schema...', new Date().toISOString());
    
    db.get("SELECT name FROM sqlite_master WHERE type='table' AND name='users'", (err, row) => {
        if (err) {
            console.error('❌ Error checking database schema:', err.message, new Date().toISOString());
            process.exit(1);
        }
        
        if (row) {
            console.log('✅ Database schema already exists. Skipping initialization.', new Date().toISOString());
            populateSampleData();
        } else {
            console.log('🏗️ Initializing database schema...', new Date().toISOString());
            const schemaSQL = fs.readFileSync(schemaPath, 'utf8');
            
            db.exec(schemaSQL, (err) => {
                if (err) {
                    console.error('❌ Error initializing schema:', err.message, new Date().toISOString());
                    process.exit(1);
                } else {
                    console.log('✅ Schema initialized successfully.', new Date().toISOString());
                    populateSampleData();
                }
            });
        }
    });
}

function generateSampleData() {
    console.log('📊 Generating realistic sample data...');
    const now = new Date().toISOString();
    const oneWeekAgo = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString();
    const twoWeeksAgo = new Date(Date.now() - 14 * 24 * 60 * 60 * 1000).toISOString();

    const users = [
        { id: uuidv4(), email: 'john.smith@example.com', username: 'johnsmith', firstName: 'John', lastName: 'Smith', createdAt: twoWeeksAgo, updatedAt: now },
        { id: uuidv4(), email: 'emily.johnson@example.com', username: 'emilyjohnson', firstName: 'Emily', lastName: 'Johnson', createdAt: twoWeeksAgo, updatedAt: now },
        { id: uuidv4(), email: 'michael.williams@example.com', username: 'michaelw', firstName: 'Michael', lastName: 'Williams', createdAt: oneWeekAgo, updatedAt: now },
        { id: uuidv4(), email: 'sarah.brown@example.com', username: 'sarahb', firstName: 'Sarah', lastName: 'Brown', createdAt: oneWeekAgo, updatedAt: now },
        { id: uuidv4(), email: 'david.jones@example.com', username: 'davidj', firstName: 'David', lastName: 'Jones', createdAt: now, updatedAt: now },
    ];

    const friends = [
        { id: uuidv4(), userId: users[0].id, friendId: users[1].id, status: 'accepted', createdAt: oneWeekAgo, updatedAt: oneWeekAgo },
        { id: uuidv4(), userId: users[0].id, friendId: users[2].id, status: 'accepted', createdAt: twoWeeksAgo, updatedAt: twoWeeksAgo },
        { id: uuidv4(), userId: users[1].id, friendId: users[3].id, status: 'pending', createdAt: now, updatedAt: now },
        { id: uuidv4(), userId: users[2].id, friendId: users[4].id, status: 'accepted', createdAt: oneWeekAgo, updatedAt: oneWeekAgo },
    ];

    const friendContexts = [
        { id: uuidv4(), userId: users[0].id, friendId: users[1].id, context: 'College roommate', createdAt: oneWeekAgo },
        { id: uuidv4(), userId: users[0].id, friendId: users[2].id, context: 'Work colleague', createdAt: twoWeeksAgo },
        { id: uuidv4(), userId: users[2].id, friendId: users[4].id, context: 'Gym buddy', createdAt: oneWeekAgo },
    ];

    const friendLikes = [
        { id: uuidv4(), userId: users[0].id, friendId: users[1].id, createdAt: now },
        { id: uuidv4(), userId: users[2].id, friendId: users[0].id, createdAt: oneWeekAgo },
        { id: uuidv4(), userId: users[1].id, friendId: users[3].id, createdAt: twoWeeksAgo },
    ];

    const collections = [
        { id: uuidv4(), userId: users[0].id, name: 'Work Projects', description: 'All current work-related files', createdAt: twoWeeksAgo, updatedAt: now },
        { id: uuidv4(), userId: users[1].id, name: 'Vacation Photos', description: 'Photos from recent trips', createdAt: oneWeekAgo, updatedAt: now },
        { id: uuidv4(), userId: users[2].id, name: 'Personal Documents', description: 'Important personal files', createdAt: twoWeeksAgo, updatedAt: now },
    ];

    const folders = [
        { id: uuidv4(), userId: users[0].id, name: 'Project X', description: 'Files for Project X', parentId: null, createdAt: twoWeeksAgo, updatedAt: now },
        { id: uuidv4(), userId: users[0].id, name: 'Client Presentations', description: 'Presentation files for clients', parentId: null, createdAt: oneWeekAgo, updatedAt: oneWeekAgo },
        { id: uuidv4(), userId: users[1].id, name: 'Italy 2023', description: 'Photos from Italy trip', parentId: null, createdAt: oneWeekAgo, updatedAt: now },
        { id: uuidv4(), userId: users[2].id, name: 'Taxes', description: 'Tax documents', parentId: null, createdAt: twoWeeksAgo, updatedAt: twoWeeksAgo },
    ];

    const files = [
        { id: uuidv4(), userId: users[0].id, folderId: folders[0].id, collectionId: collections[0].id, key: 'project_x_proposal.pdf', name: 'Project X Proposal', contentType: 'application/pdf', size: 2048576, uploadedAt: twoWeeksAgo, createdAt: twoWeeksAgo, updatedAt: twoWeeksAgo },
        { id: uuidv4(), userId: users[0].id, folderId: folders[1].id, collectionId: collections[0].id, key: 'client_presentation_q2.pptx', name: 'Q2 Client Presentation', contentType: 'application/vnd.openxmlformats-officedocument.presentationml.presentation', size: 5242880, uploadedAt: oneWeekAgo, createdAt: oneWeekAgo, updatedAt: oneWeekAgo },
        { id: uuidv4(), userId: users[1].id, folderId: folders[2].id, collectionId: collections[1].id, key: 'rome_colosseum.jpg', name: 'Colosseum, Rome', contentType: 'image/jpeg', size: 3145728, uploadedAt: oneWeekAgo, createdAt: oneWeekAgo, updatedAt: oneWeekAgo },
        { id: uuidv4(), userId: users[2].id, folderId: folders[3].id, collectionId: collections[2].id, key: 'tax_return_2023.pdf', name: 'Tax Return 2023', contentType: 'application/pdf', size: 1048576, uploadedAt: twoWeeksAgo, createdAt: twoWeeksAgo, updatedAt: twoWeeksAgo },
    ];

    const fileCategories = [
        { id: uuidv4(), name: 'Document', createdAt: twoWeeksAgo, updatedAt: twoWeeksAgo },
        { id: uuidv4(), name: 'Image', createdAt: twoWeeksAgo, updatedAt: twoWeeksAgo },
        { id: uuidv4(), name: 'Presentation', createdAt: twoWeeksAgo, updatedAt: twoWeeksAgo },
        { id: uuidv4(), name: 'Spreadsheet', createdAt: twoWeeksAgo, updatedAt: twoWeeksAgo },
    ];

    const fileCategoryAssociations = [
        { fileId: files[0].id, categoryId: fileCategories[0].id },
        { fileId: files[1].id, categoryId: fileCategories[2].id },
        { fileId: files[2].id, categoryId: fileCategories[1].id },
        { fileId: files[3].id, categoryId: fileCategories[0].id },
    ];

    const sharedFiles = [
        { id: uuidv4(), fileId: files[0].id, sharedBy: users[0].id, sharedWith: users[2].id, createdAt: oneWeekAgo },
        { id: uuidv4(), fileId: files[1].id, sharedBy: users[0].id, sharedWith: users[1].id, createdAt: now },
        { id: uuidv4(), fileId: files[2].id, sharedBy: users[1].id, sharedWith: users[3].id, createdAt: oneWeekAgo },
    ];

    const activityLog = [
        { id: uuidv4(), userId: users[0].id, actionType: 'FILE_UPLOAD', actionDetails: 'Uploaded Project X Proposal', createdAt: twoWeeksAgo },
        { id: uuidv4(), userId: users[0].id, actionType: 'FILE_SHARE', actionDetails: 'Shared Project X Proposal with Michael Williams', createdAt: oneWeekAgo },
        { id: uuidv4(), userId: users[1].id, actionType: 'FILE_UPLOAD', actionDetails: 'Uploaded Colosseum, Rome', createdAt: oneWeekAgo },
        { id: uuidv4(), userId: users[2].id, actionType: 'FILE_UPLOAD', actionDetails: 'Uploaded Tax Return 2023', createdAt: twoWeeksAgo },
        { id: uuidv4(), userId: users[0].id, actionType: 'FOLDER_CREATE', actionDetails: 'Created folder Client Presentations', createdAt: oneWeekAgo },
    ];

    return {
        users,
        friends,
        friendContexts,
        friendLikes,
        collections,
        folders,
        files,
        fileCategories,
        fileCategoryAssociations,
        sharedFiles,
        activityLog,
    };
}

function populateSampleData() {
    console.log('🏁 Starting database population process:', new Date().toISOString());

    const sampleData = generateSampleData();

    db.serialize(() => {
        console.log('🔒 Beginning transaction:', new Date().toISOString());
        db.run('BEGIN TRANSACTION');

        // Insert users
        console.log('👤 Inserting users:', new Date().toISOString());
        const insertUser = db.prepare(`
            INSERT INTO users (id, email, username, first_name, last_name, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, ?)
        `);
        sampleData.users.forEach((user, index) => {
            insertUser.run(user.id, user.email, user.username, user.firstName, user.lastName, user.createdAt, user.updatedAt, (err) => {
                if (err) {
                    console.error(`❌ Error inserting user ${index + 1}:`, err.message, new Date().toISOString());
                } else {
                    console.log(`✅ User ${index + 1} inserted successfully:`, JSON.stringify(user), new Date().toISOString());
                }
            });
        });
        insertUser.finalize();

        // Insert friends
        console.log('🤝 Inserting friends:', new Date().toISOString());
        const insertFriend = db.prepare(`
            INSERT INTO friends (id, user_id, friend_id, status, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?)
        `);
        sampleData.friends.forEach((friend, index) => {
            insertFriend.run(friend.id, friend.userId, friend.friendId, friend.status, friend.createdAt, friend.updatedAt, (err) => {
                if (err) {
                    console.error(`❌ Error inserting friend ${index + 1}:`, err.message, new Date().toISOString());
                } else {
                    console.log(`✅ Friend ${index + 1} inserted successfully:`, JSON.stringify(friend), new Date().toISOString());
                }
            });
        });
        insertFriend.finalize();

        // Insert friend contexts
        console.log('🧩 Inserting friend contexts:', new Date().toISOString());
        const insertFriendContext = db.prepare(`
            INSERT INTO friend_contexts (id, user_id, friend_id, context, created_at)
            VALUES (?, ?, ?, ?, ?)
        `);
        sampleData.friendContexts.forEach((context, index) => {
            insertFriendContext.run(context.id, context.userId, context.friendId, context.context, context.createdAt, (err) => {
                if (err) {
                    console.error(`❌ Error inserting friend context ${index + 1}:`, err.message, new Date().toISOString());
                } else {
                    console.log(`✅ Friend context ${index + 1} inserted successfully:`, JSON.stringify(context), new Date().toISOString());
                }
            });
        });
        insertFriendContext.finalize();

        // Insert friend likes
        console.log('👍 Inserting friend likes:', new Date().toISOString());
        const insertFriendLike = db.prepare(`
            INSERT INTO friend_likes (id, user_id, friend_id, created_at)
            VALUES (?, ?, ?, ?)
        `);
        sampleData.friendLikes.forEach((like, index) => {
            insertFriendLike.run(like.id, like.userId, like.friendId, like.createdAt, (err) => {
                if (err) {
                    console.error(`❌ Error inserting friend like ${index + 1}:`, err.message, new Date().toISOString());
                } else {
                    console.log(`✅ Friend like ${index + 1} inserted successfully:`, JSON.stringify(like), new Date().toISOString());
                }
            });
        });
        insertFriendLike.finalize();

        // Insert collections
        console.log('📚 Inserting collections:', new Date().toISOString());
        const insertCollection = db.prepare(`
            INSERT INTO collections (id, user_id, name, description, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?)
        `);
        sampleData.collections.forEach((collection, index) => {
            insertCollection.run(collection.id, collection.userId, collection.name, collection.description, collection.createdAt, collection.updatedAt, (err) => {
                if (err) {
                    console.error(`❌ Error inserting collection ${index + 1}:`, err.message, new Date().toISOString());
                } else {
                    console.log(`✅ Collection ${index + 1} inserted successfully:`, JSON.stringify(collection), new Date().toISOString());
                }
            });
        });
        insertCollection.finalize();

        // Insert folders
        console.log('📂 Inserting folders:', new Date().toISOString());
        const insertFolder = db.prepare(`
            INSERT INTO folders (id, user_id, name, description, parent_id, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, ?)
        `);
        sampleData.folders.forEach((folder, index) => {
            insertFolder.run(folder.id, folder.userId, folder.name, folder.description, folder.parentId, folder.createdAt, folder.updatedAt, (err) => {
                if (err) {
                    console.error(`❌ Error inserting folder ${index + 1}:`, err.message, new Date().toISOString());
                } else {
                    console.log(`✅ Folder ${index + 1} inserted successfully:`, JSON.stringify(folder), new Date().toISOString());
                }
            });
        });
        insertFolder.finalize();

        // Insert files
        console.log('📄 Inserting files:', new Date().toISOString());
        const insertFile = db.prepare(`
            INSERT INTO files (id, user_id, folder_id, collection_id, key, name, content_type, size, uploaded_at, created_at, updated_at)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        `);
        sampleData.files.forEach((file, index) => {
            insertFile.run(file.id, file.userId, file.folderId, file.collectionId, file.key, file.name, file.contentType, file.size, file.uploadedAt, file.createdAt, file.updatedAt, (err) => {
                if (err) {
                    console.error(`❌ Error inserting file ${index + 1}:`, err.message, new Date().toISOString());
                } else {
                    console.log(`✅ File ${index + 1} inserted successfully:`, JSON.stringify(file), new Date().toISOString());
                    // Update fileCategoryAssociations and sharedFiles with the inserted file ID
                    sampleData.fileCategoryAssociations[0].fileId = file.id;
                    sampleData.sharedFiles[0].fileId = file.id;
                }
            });
        });
        insertFile.finalize();

        // Insert file categories
        console.log('🏷️ Inserting file categories:', new Date().toISOString());
        const insertFileCategory = db.prepare(`
            INSERT INTO file_categories (id, name, created_at, updated_at)
            VALUES (?, ?, ?, ?)
        `);
        sampleData.fileCategories.forEach((category, index) => {
            insertFileCategory.run(category.id, category.name, category.createdAt, category.updatedAt, (err) => {
                if (err) {
                    console.error(`❌ Error inserting file category ${index + 1}:`, err.message, new Date().toISOString());
                } else {
                    console.log(`✅ File category ${index + 1} inserted successfully:`, JSON.stringify(category), new Date().toISOString());
                    if (index === 0) {
                        // Update fileCategoryAssociations with the first category ID
                        sampleData.fileCategoryAssociations[0].categoryId = category.id;
                    }
                }
            });
        });
        insertFileCategory.finalize();

        // Insert file category associations
        console.log('🔗 Inserting file category associations:', new Date().toISOString());
        const insertFileCategoryAssociation = db.prepare(`
            INSERT INTO file_category_associations (file_id, category_id)
            VALUES (?, ?)
        `);
        sampleData.fileCategoryAssociations.forEach((association, index) => {
            insertFileCategoryAssociation.run(association.fileId, association.categoryId, (err) => {
                if (err) {
                    console.error(`❌ Error inserting file category association ${index + 1}:`, err.message, new Date().toISOString());
                } else {
                    console.log(`✅ File category association ${index + 1} inserted successfully:`, JSON.stringify(association), new Date().toISOString());
                }
            });
        });
        insertFileCategoryAssociation.finalize();

        // Insert shared files
        console.log('🔄 Inserting shared files:', new Date().toISOString());
        const insertSharedFile = db.prepare(`
            INSERT INTO shared_files (id, file_id, shared_by, shared_with, created_at)
            VALUES (?, ?, ?, ?, ?)
        `);
        sampleData.sharedFiles.forEach((sharedFile, index) => {
            insertSharedFile.run(sharedFile.id, sharedFile.fileId, sharedFile.sharedBy, sharedFile.sharedWith, sharedFile.createdAt, (err) => {
                if (err) {
                    console.error(`❌ Error inserting shared file ${index + 1}:`, err.message, new Date().toISOString());
                } else {
                    console.log(`✅ Shared file ${index + 1} inserted successfully:`, JSON.stringify(sharedFile), new Date().toISOString());
                }
            });
        });
        insertSharedFile.finalize();

        // Insert activity log
        console.log('📝 Inserting activity log:', new Date().toISOString());
        const insertActivityLog = db.prepare(`
            INSERT INTO activity_log (id, user_id, action_type, action_details, created_at)
            VALUES (?, ?, ?, ?, ?)
        `);
        sampleData.activityLog.forEach((activity, index) => {
            insertActivityLog.run(activity.id, activity.userId, activity.actionType, activity.actionDetails, activity.createdAt, (err) => {
                if (err) {
                    console.error(`❌ Error inserting activity log ${index + 1}:`, err.message, new Date().toISOString());
                } else {
                    console.log(`✅ Activity log ${index + 1} inserted successfully:`, JSON.stringify(activity), new Date().toISOString());
                }
            });
        });
        insertActivityLog.finalize();

        // Commit transaction
        console.log('🔐 Committing transaction:', new Date().toISOString());
        db.run('COMMIT', (err) => {
            if (err) {
                console.error('❌ Error committing transaction:', err.message, new Date().toISOString());
                console.log('↩️ Rolling back changes:', new Date().toISOString());
                db.run('ROLLBACK');
            } else {
                console.log('✅ Transaction committed successfully.', new Date().toISOString());
                console.log('🎉 Database population complete!', new Date().toISOString());
            }
            console.log('🔌 Closing database connection:', new Date().toISOString());
            db.close((err) => {
                if (err) {
                    console.error('❌ Error closing database connection:', err.message, new Date().toISOString());
                } else {
                    console.log('👋 Database connection closed.', new Date().toISOString());
                }
            });
        });
    });
}

// Error handling for the database connection
db.on('error', (err) => {
    console.error('❌ Database error:', err.message, new Date().toISOString());
});

console.log('📜 Script execution completed:', new Date().toISOString());