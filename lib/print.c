#include <stdlib.h>
#include <stdio.h>

int _print_lisp(int argc, char **argv) {
    for (int i = argc-1; i >= 0; i--){
        printf("%s ", argv[i]);
    }
    
    printf("\n");

    return 0;
}