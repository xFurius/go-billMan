document.addEventListener('astilectron-ready', function(){
    btnAdd.addEventListener("click", () => {
            astilectron.sendMessage("add-data");
        })
        btnShow.addEventListener("click", () => {
                astilectron.sendMessage("show-data");
            })
            btnExit.addEventListener("click", () => {
                astilectron.sendMessage("exit");
            })
    })
    