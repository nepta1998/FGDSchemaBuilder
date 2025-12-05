/**
 * Note: For this to work in a browser environment, you will need a polyfill
 * for `crypto.randomUUID`. Modern browsers support it, but for broader
 * compatibility, a library like `uuid` could be used.
 * For this example, we assume a modern environment where it's available.
 */

/**
 * Parses the base(...) definitions for an entity.
 * @param {string} baseStr The string containing base class definitions.
 * @param {string} headerStr The string containing all header definitions.
 * @returns {string[]} An array of base class names.
 */
function parseBaseClasses(headerStr) {
    if (!headerStr) return [];
    const baseClassRegex = /base\(([^)]+)\)/g;
    const bases = [];
    let match;
    while ((match = baseClassRegex.exec(headerStr)) !== null) {
        const classList = match[1].split(',').map(s => s.trim()).filter(Boolean);
        bases.push(...classList);
    }
    return bases;
}

/**
 * Parses various helper properties from an entity's header.
 * @param {string} headerStr The string containing all header definitions.
 * @returns {object} An object containing parsed helper data.
 */
function parseHelpers(headerStr) {
    if (!headerStr) return {};
    const helpers = {};
    // Simple key-value helpers like size() and color()
    const simpleHelperRegex = /(size|color)\s*\(([^)]+)\)/g;
    let match;
    while ((match = simpleHelperRegex.exec(headerStr)) !== null) {
        helpers[match[1]] = match[2].trim();
    }
    // Complex model helper that can contain nested structures
    const modelMatch = headerStr.match(/model\s*\(([\s\S]*)\)/);
    if (modelMatch) {
        helpers.model = modelMatch[1].trim();
    }
    return helpers;
}

/**
 * Parses the spawnflags block.
 * @param {string} flagsLine The full line containing the spawnflags definition.
 * @returns {object[]} An array of flag option objects.
 */
function parseFlags(flagsLine) {
    const options = [];
    const blockContentMatch = flagsLine.match(/\[([\s\S]*)\]/);
    if (!blockContentMatch) return [];

    const content = blockContentMatch[1];
    // Split by lines, keep comments and flag definitions
    const lines = content.split('\n').map(line => line.trim()).filter(Boolean);

    for (const line of lines) {
        if (line.startsWith('//')) {
            options.push({ type: 'comment', text: line.replace(/^\/\/\s?/, '') });
            continue;
        }
        const flagMatch = line.match(/^(\d+)\s*:\s*"([^"]+)"\s*:\s*(\d)/);
        if (flagMatch) {
            options.push({
                type: 'flag',
                value: parseInt(flagMatch[1], 10),
                label: flagMatch[2],
                default: !!parseInt(flagMatch[3], 10),
            });
        }
    }
    return options;
}

/**
 * Parses the choices block.
 * @param {string} choiceBlock The string content of the choices block.
 * @returns {object[]} An array of choice option objects.
 */
function parseChoices(choiceBlock) {
    const options = [];
    // Regex to find: 0 : "Medieval"
    const choiceRegex = /(-?\d+)\s*:\s*"([^"]+)"/g;
    let match;
    while ((match = choiceRegex.exec(choiceBlock)) !== null) {
        const [, value, label] = match;
        options.push({
            value: parseInt(value, 10),
            label: label,
        });
    }
    return options;
}


/**
 * Parses the body of an entity definition to extract its properties.
 * @param {string} bodyText The text content between the `[` and `]` of an entity.
 * @returns {object[]} An array of property objects.
 */
function parseEntityBody(bodyText) {
    const properties = [];
    const lines = bodyText.trim().split('\n').map(line => line.trim()).filter(line => line);

    for (let i = 0; i < lines.length; i++) {
        let line = lines[i];

        // Multi-line block property (e.g. delay(choices) : "Desc" : 0 = [ ... ])
        // Updated regex to handle optional description and default value before the equals sign
        const blockPropMatch = line.match(/^([a-zA-Z0-9_]+)\s*\((\w+)\)\s*(?::\s*"([^"]*)")?\s*(?::\s*(-?[\d.\s]+|"[^"]*"))?\s*=\s*$/i);

        if (blockPropMatch && lines[i + 1] && lines[i + 1].startsWith('[')) {
            const [, name, type, displayName, defaultValue] = blockPropMatch;
            let blockContent = '';
            i++; // Move to the line with [
            while (i < lines.length) {
                if (lines[i].includes(']')) {
                    blockContent += lines[i].replace(']', '');
                    break;
                }
                blockContent += lines[i] + '\n';
                i++;
            }
            const isFlags = type.toLowerCase() === 'flags';
            const isChoices = type.toLowerCase() === 'choices';
            properties.push({
                id: crypto.randomUUID(),
                name: name || '',
                type: type || 'string',
                displayName: displayName || name,
                defaultValue: defaultValue ? defaultValue.trim() : '',
                description: '',
                flags: isFlags ? parseFlags(`[${blockContent}]`) : undefined,
                options: isChoices ? parseChoices(blockContent) : [],
            });
            continue;
        }

        // Single-line block property (e.g. delay(choices) = [ ... ])
        const singleLineBlockMatch = line.match(/^([a-zA-Z0-9_]+)\s*\((\w+)\)\s*=\s*\[([\s\S]*)\]$/i);
        if (singleLineBlockMatch) {
            const [, name, type, blockContent] = singleLineBlockMatch;
            const isFlags = type.toLowerCase() === 'flags';
            const isChoices = type.toLowerCase() === 'choices';
            properties.push({
                id: crypto.randomUUID(),
                name: name || '',
                type: type || 'string',
                displayName: name,
                defaultValue: '',
                description: '',
                flags: isFlags ? parseFlags(`[${blockContent}]`) : undefined,
                options: isChoices ? parseChoices(blockContent) : [],
            });
            continue;
        }

        // Standard, single-line property
        const propRegex = /([a-zA-Z0-9_]+)\s*\((\w+)\)\s*(?::\s*"([^"]*)")?\s*(?::\s*(-?[\d.\s]+|"[^"]*"))?\s*(?::\s*"([^"]*)")?/;
        const propMatch = line.match(propRegex);

        if (propMatch) {
            const [, name, type, displayName, defaultValue, description] = propMatch;
            properties.push({
                id: crypto.randomUUID(),
                name: name || '',
                type: type || 'string',
                displayName: displayName || name,
                defaultValue: defaultValue ? defaultValue.replace(/"/g, '').trim() : '',
                description: description || '',
            });
        }
    }
    return properties;
}


/**
 * Parses a raw FGD file string into a structured JSON object.
 * @param {string} fgdText The raw text content of an FGD file. This should be the full file content.
 * @returns {object} A structured JSON representation of the FGD data.
 */
export function parseFGD(fgdText) {
    const schema = {
        metadata: {
            mapsize: null,
            includes: [],
        },
        entities: [],
        comments: [], // <-- Add this!
    };

    // Normalize line endings
    const lines = fgdText.replace(/\r\n/g, '\n').split('\n');

    // 1. Collect global comments at the top
    let i = 0;
    while (i < lines.length && lines[i].trim().startsWith('//')) {
        // Remove the leading // and any leading space after it
        schema.comments.push(lines[i].replace(/^\/\/\s?/, ''));
        i++;
    }

    // 2. Re-join the rest for further parsing
    const restText = lines.slice(i).join('\n');

    // 3. Parse top-level directives
    const mapsizeMatch = restText.match(/@mapsize\s*\(\s*(-?\d+)\s*,\s*(-?\d+)\s*\)/);
    if (mapsizeMatch) {
        schema.metadata.mapsize = {
            min: parseInt(mapsizeMatch[1], 10),
            max: parseInt(mapsizeMatch[2], 10),
        };
    }

    const includeMatches = [...restText.matchAll(/@include\s*"([^"]+)"/g)];
    schema.metadata.includes = includeMatches.map(match => match[1]);

    // 4. Parse entity blocks (same as before)
    const entityRegex = /@([a-zA-Z]+)\s*([\s\S]*?)=\s*([^:[\n]+)\s*(?::\s*"([^"]*)")?\s*\[([\s\S]*?)\]/g;

    let match;
    while ((match = entityRegex.exec(restText)) !== null) {
        if (match.index > entityRegex.lastIndex) {
            entityRegex.lastIndex = match.index;
        }

        const [, classType, header, name, description, body] = match;

        const entity = {
            id: crypto.randomUUID(),
            classType: classType || '',
            name: name || '',
            description: (description || '').trim(),
            baseClasses: parseBaseClasses(header || ''),
            helpers: parseHelpers(header || ''),
            properties: body ? parseEntityBody(body) : [],
        };

        // Fix: Treat worldspawn as solidclass
        if (entity.name.trim().toLowerCase() === 'worldspawn') {
            entity.classType = 'SolidClass';
        }

        schema.entities.push(entity);
    }
    return schema;
}