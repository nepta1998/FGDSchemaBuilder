const regex = /^([a-zA-Z0-9_]+)\s*\((\w+)\)\s*(?::\s*"([^"]*)")?\s*(?::\s*(-?[\d.\s]+|"[^"]*"))?\s*=\s*$/i;

const lines = [
    'worldtype(choices) : "Ambience" : 0 =',
    'spawnflags(Flags) =',
    'delay(choices) : "Attenuation" =',
    'style(Choices) : "Appearance" : 0 ='
];

lines.forEach(line => {
    const match = line.match(regex);
    console.log(`Line: "${line}"`);
    if (match) {
        console.log('  MATCHED!');
        console.log('  Name:', match[1]);
        console.log('  Type:', match[2]);
        console.log('  Display:', match[3]);
        console.log('  Default:', match[4]);
    } else {
        console.log('  NO MATCH');
    }
    console.log('---');
});
