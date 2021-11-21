const slashData = {
slashCommands: [
        {
            Name: "ping",
            Description: "Pong!"
        },
        {
            Name: "pong",
            Description: "Ping!"
        },
        {
            Name: "plink",
            Description: "Plonk!"
        },
        {
            Name: "plonk",
            Description: "Plink!"
        }
    ]
};


function slashCommands() {

    return JSON.stringify(slashData)
}

console.log(slashData)

console.log(JSON.parse(JSON.stringify(slashData)))

console.log(JSON.stringify(slashData))