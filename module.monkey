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

let p = Player();
puts(p.oldVelocity);
let u = 3;
let x = while(u < 5){
    puts(u);
    let u = u+1;
}
puts("ending", x);