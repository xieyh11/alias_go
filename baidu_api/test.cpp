#include "nlp.h"
#include <iostream>
#include <string>

int main(){
    std::string app_id = "14306411";
    std::string api_key = "F9bzSsEtyeTi8c1QGUZbhBMc";
    std::string secret_key = "klAMhGy2i9xsUA4MQjGGfdlzS3zTQbEb";

    aip::Nlp client(app_id, api_key, secret_key);
    Json::Value result;
    result = client.word_embedding("世纪", aip::null);
    std::cout << result << std::endl;
    // auto vec = result["vec"];
    // std::cout << vec.size() << std::endl;
    // for(int i = 0; i < vec.size(); i++){
    //     std::cout << vec[i] << " ";
    // }
    return 0;
}