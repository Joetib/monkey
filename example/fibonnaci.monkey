let fibonnaci = fn (limit) {
    let a = 0;
    let b = 1;
    while (a < limit) {
        puts(a)
        let temp = a
        let a = b
        let b = temp + a
    }
}

fibonnaci(15)