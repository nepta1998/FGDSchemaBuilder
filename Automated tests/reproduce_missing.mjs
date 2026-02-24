import { readFileSync } from 'fs';
import { parseFGD } from '../src/core/FGDParser.js';
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
    const schema = parseFGD(fgdContent);

    // Check for worldspawn and worldtype
    const worldspawn = schema.entities.find(e => e.name === 'worldspawn');
    if (worldspawn) {
        console.log('Found worldspawn');
        const worldtype = worldspawn.properties.find(p => p.name === 'worldtype');
        if (worldtype) {
            console.log('Found worldtype property');
            console.log('worldtype type:', worldtype.type);
            console.log('worldtype options:', worldtype.options ? worldtype.options.length : 'undefined');
        } else {
            console.error('MISSING: worldtype property in worldspawn');
        }
    } else {
        console.error('MISSING: worldspawn entity');
    }

    // Check for item_health and spawnflags
    const itemHealth = schema.entities.find(e => e.name === 'item_health');
    if (itemHealth) {
        console.log('Found item_health');
        const spawnflags = itemHealth.properties.find(p => p.name === 'spawnflags');
        if (spawnflags) {
            console.log('Found spawnflags property');
            console.log('spawnflags flags:', spawnflags.flags ? spawnflags.flags.length : 'undefined');
        } else {
            console.error('MISSING: spawnflags property in item_health');
        }
    } else {
        console.error('MISSING: item_health entity');
    }

} catch (error) {
    console.error('Error parsing FGD:', error);
}
