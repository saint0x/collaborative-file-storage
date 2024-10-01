const axios = require('axios');
const fs = require('fs').promises;
const path = require('path');
const { exec } = require('child_process');
const FormData = require('form-data');

const SERVER_URL = 'http://localhost:8080';
const GO_SERVER_PATH = './backend/cmd/server/main.go';

async function startGoServer() {
    return new Promise((resolve, reject) => {
        const server = exec(`go run ${GO_SERVER_PATH}`, (error, stdout, stderr) => {
            if (error) {
                console.error(`Error starting Go server: ${error}`);
                reject(error);
            }
            console.log(`Go server output: ${stdout}`);
            console.error(`Go server errors: ${stderr}`);
        });

        server.stdout.on('data', (data) => {
            if (data.includes('Starting server on')) {
                console.log('Go server started successfully');
                resolve(server);
            }
        });
    });
}

async function waitForServerReady() {
    for (let i = 0; i < 30; i++) {
        try {
            await axios.get(`${SERVER_URL}/health`);
            console.log('Server is ready');
            return;
        } catch (error) {
            await new Promise(resolve => setTimeout(resolve, 1000));
        }
    }
    throw new Error('Server failed to start within 30 seconds');
}

async function uploadFile(filePath, fileName) {
    const fileContent = await fs.readFile(filePath);
    const formData = new FormData();
    formData.append('file', fileContent, fileName);

    try {
        const response = await axios.post(`${SERVER_URL}/upload`, formData, {
            headers: {
                ...formData.getHeaders(),
                'Authorization': 'Bearer dummy_token' // Replace with actual auth token if needed
            },
        });
        console.log(`✅ File uploaded successfully: ${fileName}`);
        return response.data;
    } catch (error) {
        console.error(`❌ Failed to upload file: ${fileName}`);
        console.error('Error details:', error.message);
        if (error.response) {
            console.error('Response status:', error.response.status);
            console.error('Response data:', error.response.data);
        }
        return null;
    }
}

async function createUser(userData) {
    try {
        const response = await axios.post(`${SERVER_URL}/users`, userData, {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer dummy_token' // Replace with actual auth token if needed
            }
        });
        console.log(`✅ User created successfully: ${userData.username}`);
        return response.data;
    } catch (error) {
        console.error(`❌ Failed to create user: ${userData.username}`);
        console.error('Error details:', error.message);
        if (error.response) {
            console.error('Response status:', error.response.status);
            console.error('Response data:', error.response.data);
        }
        return null;
    }
}

async function populateDatabase() {
    console.log('🚀 Starting database population...');

    // Create users
    const users = [
        { username: 'john_doe', email: 'john@example.com', firstName: 'John', lastName: 'Doe' },
        { username: 'jane_smith', email: 'jane@example.com', firstName: 'Jane', lastName: 'Smith' },
    ];

    for (const userData of users) {
        await createUser(userData);
    }

    // Upload sample files
    const sampleFiles = [
        { path: './sample_files/project_x_proposal.pdf', name: 'Project X Proposal.pdf' },
        { path: './sample_files/client_presentation_q2.pptx', name: 'Q2 Client Presentation.pptx' },
        { path: './sample_files/rome_colosseum.jpg', name: 'Rome Colosseum.jpg' },
        { path: './sample_files/tax_return_2023.pdf', name: 'Tax Return 2023.pdf' },
    ];

    for (const file of sampleFiles) {
        await uploadFile(file.path, file.name);
    }

    console.log('✅ Database population complete!');
}

async function main() {
    try {
        console.log('Starting Go server...');
        const server = await startGoServer();
        await waitForServerReady();
        
        await populateDatabase();
        
        console.log('Shutting down Go server...');
        server.kill();
    } catch (error) {
        console.error('❌ Error executing script:', error.message);
    }
}

main();