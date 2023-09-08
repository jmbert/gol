#include <stdlib.h>
#include <stdio.h>

char *_add_lisp(int argc, char **argv) {
    int result = 0;
    for (int i = argc-1; i >= 0; i--) {
        result += atoi(argv[i]);
    }
    int size = snprintf(NULL, 0, "%d", result);
    char *out = (char*)malloc(sizeof(char)*size);
    sprintf(out, "%d", result);
    return out;
}

char *_sub_lisp(int argc, char **argv) {
    int result = 0;
    if (argc > 1) {
        result = atoi(argv[argc-1]);
    }

    for (int i = argc-2; i > 1; i--) {
        result -= atoi(argv[i]);
    }
    int size = snprintf(NULL, 0, "%d", result);
    char *out = (char*)malloc(sizeof(char)*size);
    sprintf(out, "%d", result);
    return out;
}

char *_mul_lisp(int argc, char **argv) {
    int result = 0;
    for (int i = argc-1; i >= 0; i--) {
        result *= atoi(argv[i]);
    }
    int size = snprintf(NULL, 0, "%d", result);
    char *out = (char*)malloc(sizeof(char)*size);
    sprintf(out, "%d", result);
    return out;
}

char *_div_lisp(int argc, char **argv) {
    int result = 0;
    result = atoi(argv[argc-1]);
    

    for (int i = argc-1; i >= 0; i--) {
        result /= atoi(argv[i]);
    }
    int size = snprintf(NULL, 0, "%d", result);
    char *out = (char*)malloc(sizeof(char)*size);
    sprintf(out, "%d", result);
    return out;
}