const AWS = require("aws-sdk");
const fs = require("fs");
const path = require("path");
const mime = require('mime');
const randomstring = require('randomstring');
const rootFolderName = process.env.BUILD_DIRECTORY || 'dist'

const config = {
    s3BucketName: process.env.BUCKET_NAME,
    folderPath: `./${rootFolderName}`
};

const awsConfig = {
    signatureVersion: 'v4',
    accessKeyId: process.env.AWS_ACCESS_KEY_ID,
    secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY,
    region: process.env.REGION,
}

const s3 = new AWS.S3(awsConfig);
const distFolderPath = path.join(__dirname, config.folderPath);
uploadDirectoryFiles(distFolderPath)

function uploadDirectoryFiles(distFolderPath) {
    const files = fs.readdirSync(distFolderPath)
    if (!files || files.length === 0) {
        console.log(`provided folder '${distFolderPath}' is empty or does not exist.`);
        return;
    }
    for (const fileName of files) {
        const filePath = path.join(distFolderPath, fileName);
        if (fs.lstatSync(filePath).isDirectory()) {
            uploadDirectoryFiles(filePath)
            continue;
        }
        uploadFile(filePath, fileName)
    }
}

function uploadFile(filePath, fileName) {
    const relativeFilePath = `${__dirname}/${rootFolderName}/`
    const fileKey = filePath.replace(relativeFilePath, '')
    console.log({ fileName, filePath, fileKey })
    const fileContent = fs.readFileSync(filePath)
    const ContentType = mime.getType(filePath)
    s3.putObject({
        Bucket: config.s3BucketName,
        Key: fileKey,
        Body: fileContent,
        ContentType
    }, (err, res) => {
        if (err) {
            return console.log("Error uploading file ", err)
        }
        console.log(`Successfully uploaded '${fileKey}'!`, {res});
    });
}


const reference = randomstring.generate(16);
const cloudfront = new AWS.CloudFront(awsConfig);
const params = {
    DistributionId: process.env.DISTRIBUTION_ID,
    InvalidationBatch: {
        CallerReference: reference,
        Paths: {
            Quantity: 1,
            Items: [
                '/*',
            ]
        }
    }
};

cloudfront.createInvalidation (params, function(err, data) {
    if (err) console.log(err, err.stack);
    else console.log(data);
});
