
/*
    get block of options which start and end by '(' ')' characters
*/
function getOptionBlock(str) {
    let level = 0;
    for (let i = 0; i < str.length; i++) {
        if (str[i] == '(') {
            level++;
        } else if (str[i] == ')') {
            if (level == 0) {
                return str.substring(0, i);
            } else {
                level--;
            }
        }
    }
    return "";
}
/*
    handle AND, OR condition
*/
function handleConditions(isAnd, index, str, objResult, currentCondition, currentOption) {
    let condition = "";
    if (isAnd) {
        condition = str.substring(index, index + 3);
        if (condition !== "AND") {
            return "";
        }
    } else {
        condition = str.substring(index, index + 2);
        if (condition !== "OR") {
            return "";
        }
    }

    // Create option array
    if (currentCondition === "") {
        currentCondition = condition;
        if (currentOption != null) {
            objResult[condition] = [currentOption];
            return condition;
        }
        else {
            return ""; // wrong format
        }
    }
    // don't allow 2 conditions
    else if (currentCondition != condition) {
        return "";
    }
    else {
        return condition
    }
}

/*
    get object result (main method)
*/
function getObjResult(str) {
    let objResult = {};
    let index = 0;
    let currentOption = null;
    let currentCondition = "";
    while (index < str.length) {
        switch (str[index]) {
            // begin option
            case '{':
                const optionTemp = str.substring(index + 1, index + 5);
                index += 6;
                // Validate
                if(isNaN(parseInt(optionTemp))){
                    return null;
                }
                // Add option
                if (currentCondition !== "") {
                    objResult[Object.keys(objResult)[0]].push(optionTemp);
                }
                else {
                    // Add to objResult later;
                    currentOption = optionTemp;
                }
                break;
            case 'A':
                const andCondition = handleConditions(true, index, str, objResult, currentCondition, currentOption);
                if (andCondition === "") {
                    return null;
                }
                else
                {
                    currentCondition = andCondition;
                }
                index += andCondition.length;
                break;
            case 'O':
                const orCondition = handleConditions(false, index, str, objResult, currentCondition, currentOption);
                if (orCondition === "") {
                    return null;
                }
                else
                {
                    currentCondition = orCondition;
                }
                index += orCondition.length;
                break;
            case '(': // option block
                // get block in string
                const block = getOptionBlock(str.substring(index + 1, str.length));
                if (block === "") {
                    return null;
                }
                //recursive call
                let objTemp = getObjResult(block);
                if (objTemp != null) {
                    // Add option
                    if (currentCondition !== "") {
                        objResult[Object.keys(objResult)[0]].push(objTemp);
                    }
                    else {
                        // Add to objResult later;
                        currentOption = objTemp;
                    }
                    index += block.length + 2; // 2 is ()
                }
                else {
                    return null;
                }
                break;

            default: // Wrong format
                return null;
        }
    }
    return objResult;
}

const optionRule = '{1069} AND ({1070} OR {1071} OR {1072}) AND {1244} AND ({1245} OR {1339})';
// remove all space characters. In order to handle { 1069} case. In that case, we might get " 1069"
let trimText = optionRule.replace(/ /g, '');

let output = getObjResult(trimText);

console.log(output);
//Show children element in array
//console.log(output['AND'][1]);

