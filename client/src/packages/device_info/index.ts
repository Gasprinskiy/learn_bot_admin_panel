import { DeviceOSMap, DeviceTypeMap } from './constants';
import type { DeviceInfo } from './types';

function defineDeviceOS(userAgent: string): string {
  for (const key in DeviceOSMap) {
    if (userAgent.includes(key)) {
      return DeviceOSMap[key];
    }
  }

  return 'unkown';
}

function defineDeviceType(userAgent: string) {
  for (const key in DeviceTypeMap) {
    if (userAgent.includes(key)) {
      return DeviceOSMap[key];
    }
  }

  return 'Desktop';
}

function getDeviceInfo(): DeviceInfo {
  const userAgent = navigator.userAgent;
  const maxTouchPoints = navigator.maxTouchPoints || 0;

  return {
    os: defineDeviceOS(userAgent),
    deviceType: defineDeviceType(userAgent),
    maxTouchPoints,
    userAgent,
  };
}

function simpleHash(str: string): string {
  let hash = 2166136261;
  for (let i = 0; i < str.length; i++) {
    hash ^= str.charCodeAt(i);
    hash += (hash << 1) + (hash << 4) + (hash << 7) + (hash << 8) + (hash << 24);
  }
  return (hash >>> 0).toString(16).padStart(8, '0');
}

export function generateDeviceID(): string {
  const info = getDeviceInfo();
  const data = JSON.stringify(info);
  return simpleHash(`${data}suck`);
  // simpleHash(data);
};

// simpleHash(`${data}suck`);
