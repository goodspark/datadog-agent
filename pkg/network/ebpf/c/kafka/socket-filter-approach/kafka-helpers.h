#ifndef __KAFKA_HELPERS_H
#define __KAFKA_HELPERS_H

#include "kafka-types.h"

static __inline int32_t read_big_endian_int32(const char* buf) {
//  const int32_t length = *((int32_t*)buf);
//  return bpf_ntohl(length);
    int32_t val;
    bpf_probe_read_kernel(&val, 4, (void*)buf);
    return bpf_ntohs(val);
}

//#define MINIMUM_API_VERSION_FOR_CLIENT_ID 1
//#define MAX_LENGTH_FOR_CLIENT_ID_STRING 50

    static __inline int16_t read_big_endian_int16(const char* buf) {
    //    const int16_t length = *((int16_t*)buf);
    //    return bpf_ntohs(length);
        int16_t val;
        // De-Referencing buf causes misalignment verifier error for unknown reason
        bpf_probe_read_kernel(&val, 2, (void*)buf);
        return bpf_ntohs(val);
    }

// Checking if the buffer represents kafka message
//static __always_inline bool is_kafka(const char* buf, __u32 buf_size) {
static __always_inline bool try_parse_request_header(kafka_transaction_t *kafka_transaction) {
    char *request_fragment = kafka_transaction->request_fragment;
//    const uint32_t request_fragment_size = sizeof(kafka_transaction->request_fragment);
//    if (buf_size < KAFKA_MIN_SIZE) {
//        log_debug("buffer size is less than KAFKA_MIN_SIZE");
//        return false;
//    }

    if (request_fragment == NULL) {
        log_debug("request_fragment == NULL");
        return false;
    }

    // Kafka size field is 4 bytes
//    const int32_t message_size = read_big_endian_int32(buf) + 4;
    const int32_t message_size = read_big_endian_int32(request_fragment);
    //log_debug("message_size = %d", message_size);
    //log_debug("buf_size = %d", buf_size);

    // Enforcing count to be exactly message_size + 4 to mitigate mis-classification.
    // However, this will miss long messages broken into multiple reads.
//    if (message_size < 0 || buf_size != (__u32)message_size) {
    if (message_size <= 0) {
//        log_debug("message_size < 0 || buf_size != (__u32)message_size");
        //log_debug("message_size <= 0");
        return false;
    }

    const int16_t request_api_key = read_big_endian_int16(request_fragment + 4);
    log_debug("request_api_key: %d", request_api_key);
    if (request_api_key < 0 || request_api_key > KAFKA_MAX_VERSION) {
        log_debug("request_api_key < 0 || request_api_key > KAFKA_MAX_VERSION");
        return false;
    }
    kafka_transaction->request_api_key = request_api_key;

    const int16_t request_api_version = read_big_endian_int16(request_fragment + 6);
    log_debug("request_api_version: %d", request_api_version);
    if (request_api_version < 0 || request_api_version > KAFKA_MAX_API) {
        log_debug("request_api_version < 0 || request_api_version > KAFKA_MAX_API");
        return false;
    }
    kafka_transaction->request_api_version = request_api_version;

    const int32_t correlation_id = read_big_endian_int32(request_fragment + 8);
    log_debug("correlation_id: %d", correlation_id);
    if (correlation_id < 0) {
        log_debug("correlation_id < 0");
        return false;
    }
     kafka_transaction->correlation_id = correlation_id;

    const int16_t MINIMUM_API_VERSION_FOR_CLIENT_ID = 1;
//    const uint32_t MAX_LENGTH_FOR_CLIENT_ID_STRING = 50;
    //char client_id[MAX_LENGTH_FOR_CLIENT_ID_STRING] = {0};
    __builtin_memset(kafka_transaction->client_id, 0, sizeof(kafka_transaction->client_id));
    uint16_t client_id_size_final = 0;
    if (request_api_version >= MINIMUM_API_VERSION_FOR_CLIENT_ID) {
        const int16_t client_id_size = read_big_endian_int16(request_fragment + 12);
        log_debug("client_id_size: %d", client_id_size);
        uint32_t max_size_of_client_id_string = sizeof(kafka_transaction->client_id);
        // A nullable string length might be -1 to signify null, it should be supported here
        if (client_id_size <= 0 || client_id_size > max_size_of_client_id_string) {
            log_debug("client_id <=0 || client_id_size > MAX_LENGTH_FOR_CLIENT_ID_STRING");
        } else {
            const char* client_id_in_buf = request_fragment + 14;
            client_id_size_final = client_id_size < max_size_of_client_id_string ? client_id_size : max_size_of_client_id_string;
            kafka_transaction->current_offset_in_request_fragment += client_id_size_final;
            bpf_probe_read_kernel_with_telemetry(kafka_transaction->client_id, client_id_size_final, (void*)client_id_in_buf);
            log_debug("client_id: %s", kafka_transaction->client_id);
        }
    }

    // Setting the offset for the next function to know where the actual request starts
    // TODO: should be done in a more clean way
    kafka_transaction->current_offset_in_request_fragment += 14;

    // TODO: need to check what is TAG_BUFFER that can appear in a v2 request header

    return true;
}

static __always_inline bool try_parse_request(kafka_transaction_t *kafka_transaction) {
    char *request_fragment = (char*)kafka_transaction->request_fragment;
    if (request_fragment == NULL) {
        return false;
    }

    // TODO: Support all protocol versions, currently supporting only version 7 for the testing env
    log_debug("current_offset: %d", kafka_transaction->current_offset_in_request_fragment);
    if (kafka_transaction->current_offset_in_request_fragment > sizeof(kafka_transaction->request_fragment)) {
        return false;
    }
//    log_debug("current_bytes: %d", request_fragment[kafka_transaction->current_offset_in_request_fragment]);
//    log_debug("current_bytes + 1:s %d", request_fragment[kafka_transaction->current_offset_in_request_fragment + 1]);
//    log_debug("current_bytes + 2: %d", request_fragment[kafka_transaction->current_offset_in_request_fragment + 2]);
//    (void)request_fragment;
    int16_t transactional_id_size = read_big_endian_int16(request_fragment + kafka_transaction->current_offset_in_request_fragment);
//    log_debug("transactional_id_size: %d", transactional_id_size);
    kafka_transaction->current_offset_in_request_fragment += 2;
    if (transactional_id_size > 0) {
        kafka_transaction->current_offset_in_request_fragment += transactional_id_size;
    }

    // Skipping the acks field as we have no interest in it
    kafka_transaction->current_offset_in_request_fragment += 2;

    // Skipping the timeout_ms field as we have no interest in it
    kafka_transaction->current_offset_in_request_fragment += 4;

    // Skipping number of entries for now
    kafka_transaction->current_offset_in_request_fragment += 4;

    if (kafka_transaction->current_offset_in_request_fragment > sizeof(kafka_transaction->request_fragment)) {
        return false;
    }

    // TODO: Taking only the first topic for now
    const int16_t topic_name_size = read_big_endian_int16(request_fragment + kafka_transaction->current_offset_in_request_fragment);
//    log_debug("topic_name_size: %d", topic_name_size);
        if (topic_name_size <= 0) {
        log_debug("topic_name_size <= 0");
        return false;
    }
    kafka_transaction->current_offset_in_request_fragment += 2;
    __builtin_memset(kafka_transaction->topic_name, 0, sizeof(kafka_transaction->topic_name));
//    if (kafka_transaction->current_offset_in_request_fragment > sizeof(kafka_transaction->request_fragment)) {
//            return false;
//    }
    if (topic_name_size > sizeof(kafka_transaction->topic_name)) {
        return false;
    }
    if (kafka_transaction->current_offset_in_request_fragment > sizeof(kafka_transaction->request_fragment)) {
        return false;
    }
    bpf_probe_read_kernel_with_telemetry(
        kafka_transaction->topic_name,
        topic_name_size,
        (void*)(request_fragment + kafka_transaction->current_offset_in_request_fragment));
    log_debug("topic_name: %s", request_fragment + kafka_transaction->current_offset_in_request_fragment);

    return true;
}

#endif