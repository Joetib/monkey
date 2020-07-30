puts("Module 1 running ...")

let print = fn(arg) {
    puts(arg);
}

let x = "Kofi is a boy";

class Player(){
    let oldVelocity = 0
    let play = fn (){
        print("Player is playing the ball");
    }
    let speedUp = fn (newVelocity) {
        let self.oldVelocity = self.oldVelocity + newVelocity;
        puts("Speeding up to "+ str(self.oldVelocity));
    }
}