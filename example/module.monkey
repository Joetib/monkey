puts("Module 1 running ...")

let print = fn(arg) {
    puts(arg);
}

let x = "Kofi is a boy";

class Player(){
    let play = fn (){
        print("Player is playing the ball");
    }
    let speedUp = fn (newVelocity) {
        print("Speeding up to ", newVelocity);
    }
}