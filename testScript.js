
function calculate(variable){
    variable.name = "def";
}

function getObjResult() {

    let test = 2;
    let obj = {
        name : "abc"
    }
    calculate(obj);
    console.log(obj);
}
getObjResult();

// const regex = /(\{[0-9]+\})/g;
// let match = regex.exec(result);
// while(match){
//     const filename = match;
//     console.log(filename);
//     match = regex.exec(result);
// }





// const string = 'filea.mp3 file_01.mp3 file_02.mp3'
// const regex = /(\w+)\.mp3/g;
// let match = regex.exec(string);

// while(match){
//     const filename = match[1];
//     console.log(filename);
//     match = regex.exec(string);
// }

