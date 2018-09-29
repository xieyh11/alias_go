#include "nlp"

extern "C" {
    void aip_client(void *client);
    void word_vector(const void *client, const char * word, const size_t word_size, double *vector, size_t vector_size);
}