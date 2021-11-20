var slashData = [
    {name: "Pong", description: "Ping!"},
    {name: "Ping", description: "Pong!"},
    {name: "Plink", description: "Plonk!"},
    {name: "Plonk!", description: "Plink!"}
 ];

function slashCommands() {

    var slashJSON = JSON.stringify(slashData)
    return slashJSON
}