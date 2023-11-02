// Code generated by colf(1); DO NOT EDIT.
// The compiler used schema file scheme.colf for package bench.

#include "poc.h"
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#ifndef COLFER_ENDIAN_CHECK
#if (defined(__LITTLE_ENDIAN__) && !defined(__BIG_ENDIAN__)) ||                \
    (defined(_LITTLE_ENDIAN) && !defined(_BIG_ENDIAN)) ||                      \
    defined(__AARCH64EL__) || defined(__ARMEL__) || defined(__e2k__) ||        \
    defined(__loongarch__) || defined(_MIPSEL) || defined(__MIPSEL) ||         \
    defined(__MIPSEL__) || defined(__riscv) || defined(__THUMBEL__)
#define COLFER_ENDIAN_CHECK 1
#else
#error Colfer implementation requires byte order in little endian
#endif
#endif

const uint_fast64_t COLFER_MASKS[9] = {
    0,
    0xff,
    0xffff,
    0xffffff,
    0xffffffff,
    0xffffffffff,
    0xffffffffffff,
    0xffffffffffffff,
    0xffffffffffffffff,
};

size_t colfer_marshal(const colfer *o, void *start) {
  // words of fixed section
  uint_fast64_t word0 = 22 - 1;
  uint_fast64_t word1 = 0;
  uint_fast64_t word2 = 0;
  uint_fast64_t word3 = 0;

  // write cursor at variable section
  uint8_t *p = start + 25;

  // computation register
  uint_fast64_t v;

  // pack Key int64 (with zig-zag encoding)
  v = (o->key >> 63) ^ (o->key << 1);
  if (v < 128) {
    v = v << 1 | 1;
  } else {
    p[0] = v;
    p[1] = v >> 8;
    p[2] = v >> 16;
    p[3] = v >> 24;
    p[4] = v >> 32;
    p[5] = v >> 40;
    p[6] = v >> 48;
    p[7] = v >> 56;

    size_t bit_count = __builtin_ffsll(v);
    size_t extraN = (((bit_count - 1) >> 3) + bit_count) >> 3;
    p += extraN;
    v >>= (extraN << 3) - 1;
    v = (v | 1) << extraN;
  }
  word0 |= v << 24;

  // pack Host text size
  v = o->host.len;
  if (v < 128) {
    v = v << 1 | 1;
  } else {
    p[0] = v;
    p[1] = v >> 8;
    p[2] = v >> 16;
    p[3] = v >> 24;
    p[4] = v >> 32;
    p[5] = v >> 40;
    p[6] = v >> 48;
    p[7] = v >> 56;

    size_t bit_count = __builtin_ffsll(v);
    size_t extraN = (((bit_count - 1) >> 3) + bit_count) >> 3;
    p += extraN;
    v >>= (extraN << 3) - 1;
    v = (v | 1) << extraN;
  }
  word0 |= v << 32;

  // pack Port uint16
  word0 |= (uint_fast64_t)(o->port) << 40;

  // pack Size int64 (with zig-zag encoding)
  v = (o->size >> 63) ^ (o->size << 1);
  if (v < 128) {
    v = v << 1 | 1;
  } else {
    p[0] = v;
    p[1] = v >> 8;
    p[2] = v >> 16;
    p[3] = v >> 24;
    p[4] = v >> 32;
    p[5] = v >> 40;
    p[6] = v >> 48;
    p[7] = v >> 56;

    size_t bit_count = __builtin_ffsll(v);
    size_t extraN = (((bit_count - 1) >> 3) + bit_count) >> 3;
    p += extraN;
    v >>= (extraN << 3) - 1;
    v = (v | 1) << extraN;
  }
  word0 |= v << 56;

  // pack Hash opaque64
  word1 = o->hash;

  // pack Ratio float64
  memcpy(&word2, &o->ratio, 8);

  // pack booleans
  word3 = (uint_fast64_t)(o->bools & 0xff) << 0;

  // copy payloads
  if (o->host.len > COLFER_MAX - ((void *)p - start))
    return 0;
  memcpy(p, o->host.utf8, o->host.len);
  p += o->host.len;

  size_t size = (void *)p - start;

  // finish header
  word0 |= (uint_fast64_t)size << 17 | 1 << 16;
  memcpy((uint8_t *)start + 0, &word0, 8);
  memcpy((uint8_t *)start + 8, &word1, 8);
  memcpy((uint8_t *)start + 16, &word2, 8);
  ((uint8_t *)start)[24] = word3;

  return size;
}

size_t colfer_unmarshal(colfer *o, const void *start) {
  // words of fixed section
  uint_fast64_t word0;
  uint_fast64_t word1;
  uint_fast64_t word2;
  uint_fast64_t word3;
  memcpy(&word0, start + 0, 8);
  memcpy(&word1, start + 8, 8);
  memcpy(&word2, start + 16, 8);
  memcpy(&word3, start + 24, 8);

  // read cursor at variable section
  uint8_t *p = (uint8_t *)start + (word0 & 0xffff) + 4;

  // unpack variable size
  uint_fast64_t v = word0 >> 17 & 0x7f;
  if ((word0 & (uint_fast64_t)1 << 16) == 0) {
    uint_fast64_t tz = __builtin_ctz(v | 0x80) + 1;
    v <<= (tz << 3) - tz;
    v &= ~COLFER_MASKS[tz];
    uint_fast64_t tail;
    memcpy(&tail, p, 8);
    v |= tail & COLFER_MASKS[tz];
    p += tz;
  }
  if (v > COLFER_MAX)
    return 0;
  size_t size = v;

  // unpack Key int64
  v = word0 >> 25 & 0x7f;
  if ((word0 & (uint_fast64_t)1 << 24) == 0) {
    uint_fast64_t tz = __builtin_ctz(v | 0x80) + 1;
    v <<= (tz << 3) - tz;
    v &= ~COLFER_MASKS[tz];
    uint_fast64_t tail;
    memcpy(&tail, p, 8);
    v |= tail & COLFER_MASKS[tz];
    p += tz;
  }
  o->key = (int64_t)(v >> 1) ^ -(int64_t)(v & 1);

  // unpack Host text size
  v = word0 >> 33 & 0x7f;
  if ((word0 & (uint_fast64_t)1 << 32) == 0) {
    uint_fast64_t tz = __builtin_ctz(v | 0x80) + 1;
    v <<= (tz << 3) - tz;
    v &= ~COLFER_MASKS[tz];
    uint_fast64_t tail;
    memcpy(&tail, p, 8);
    v |= tail & COLFER_MASKS[tz];
    p += tz;
  }
  o->host.len = v;

  // unpack Port uint16
  o->port = word0 >> 40;

  // unpack Size int64
  v = word0 >> 57 & 0x7f;
  if ((word0 & (uint_fast64_t)1 << 56) == 0) {
    uint_fast64_t tz = __builtin_ctz(v | 0x80) + 1;
    v <<= (tz << 3) - tz;
    v &= ~COLFER_MASKS[tz];
    uint_fast64_t tail;
    memcpy(&tail, p, 8);
    v |= tail & COLFER_MASKS[tz];
    p += tz;
  }
  o->size = (int64_t)(v >> 1) ^ -(int64_t)(v & 1);

  // unpack Hash opaque64
  o->hash = word1;

  // unpack Ratio float64
  memcpy(&o->ratio, &word2, 8);

  // unpack booleans
  o->bools = (word3 & 0xff) << 0;

  // clear/undo absent fields
  if ((word0 & 0xffff) < 22 - 1) {
    switch (word0 & 0xffff) {
    default:
      return 0;
    case 1 - 1:
      o->host.len = 0;
    case 2 - 1:
      o->port = 0;
    case 4 - 1:
      o->size = 0;
    case 5 - 1:
      o->hash = 0;
    case 13 - 1:
      o->ratio = 0;
    case 21 - 1:
      o->bools = 0;
    }
  };

  // copy payloads
  uint8_t *offset = (uint8_t *)start + size - o->host.len;
  if (offset < p)
    return 0;

  {
    char *s = malloc(o->host.len + 1);
    memcpy(s, offset, o->host.len);
    s[o->host.len] = 0; // null terminator
    o->host.utf8 = s;
  }

  return size;
}