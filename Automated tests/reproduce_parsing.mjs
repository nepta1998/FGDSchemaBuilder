import { readFileSync } from 'fs';
import { parseFGD } from './temp_FGDParser.mjs';
import crypto from 'crypto';

// Polyfill crypto.randomUUID for Node.js environment
if (!global.crypto) {
    global.crypto = crypto;
}
if (!global.crypto.randomUUID) {
    global.crypto.randomUUID = () => crypto.randomUUID();
}

try {
    const fgdContent = readFileSync('./Quake.fgd', 'utf-8');
    console.log('Read Quake.fgd, length:', fgdContent.length);
    const schema = parseFGD(fgdContent);
    console.log('Successfully parsed FGD.');
    console.log('Entities found:', schema.entities.length);
} catch (error) {
    console.error('Error parsing FGD:', error);
}
