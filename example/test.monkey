
let fact = fn(end){
    puts("Now at: ", end)
    if (end <2) {
        return 1;
    } else {
        return  end * fact(end - 1);
        

    }
}
puts(fact(1234))