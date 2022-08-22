// package: loggo
// file: loggo/types.proto

import * as jspb from "google-protobuf";
import * as google_protobuf_struct_pb from "google-protobuf/google/protobuf/struct_pb";

export class LogEntryResponse extends jspb.Message {
  getPosition(): number;
  setPosition(value: number): void;

  hasEntry(): boolean;
  clearEntry(): void;
  getEntry(): google_protobuf_struct_pb.Struct | undefined;
  setEntry(value?: google_protobuf_struct_pb.Struct): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LogEntryResponse.AsObject;
  static toObject(includeInstance: boolean, msg: LogEntryResponse): LogEntryResponse.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LogEntryResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LogEntryResponse;
  static deserializeBinaryFromReader(message: LogEntryResponse, reader: jspb.BinaryReader): LogEntryResponse;
}

export namespace LogEntryResponse {
  export type AsObject = {
    position: number,
    entry?: google_protobuf_struct_pb.Struct.AsObject,
  }
}

export class LogEntryRequest extends jspb.Message {
  hasFromPosition(): boolean;
  clearFromPosition(): void;
  getFromPosition(): number;
  setFromPosition(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LogEntryRequest.AsObject;
  static toObject(includeInstance: boolean, msg: LogEntryRequest): LogEntryRequest.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: LogEntryRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LogEntryRequest;
  static deserializeBinaryFromReader(message: LogEntryRequest, reader: jspb.BinaryReader): LogEntryRequest;
}

export namespace LogEntryRequest {
  export type AsObject = {
    fromPosition: number,
  }
}

export class Template extends jspb.Message {
  hasName(): boolean;
  clearName(): void;
  getName(): string;
  setName(value: string): void;

  clearKeysList(): void;
  getKeysList(): Array<Key>;
  setKeysList(value: Array<Key>): void;
  addKeys(value?: Key, index?: number): Key;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Template.AsObject;
  static toObject(includeInstance: boolean, msg: Template): Template.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Template, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Template;
  static deserializeBinaryFromReader(message: Template, reader: jspb.BinaryReader): Template;
}

export namespace Template {
  export type AsObject = {
    name: string,
    keysList: Array<Key.AsObject>,
  }
}

export class Color extends jspb.Message {
  getForeground(): string;
  setForeground(value: string): void;

  getBackground(): string;
  setBackground(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Color.AsObject;
  static toObject(includeInstance: boolean, msg: Color): Color.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Color, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Color;
  static deserializeBinaryFromReader(message: Color, reader: jspb.BinaryReader): Color;
}

export namespace Color {
  export type AsObject = {
    foreground: string,
    background: string,
  }
}

export class Key extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getType(): string;
  setType(value: string): void;

  hasLayout(): boolean;
  clearLayout(): void;
  getLayout(): string;
  setLayout(value: string): void;

  hasColor(): boolean;
  clearColor(): void;
  getColor(): Color | undefined;
  setColor(value?: Color): void;

  getMaxWidth(): number;
  setMaxWidth(value: number): void;

  clearColorWhenList(): void;
  getColorWhenList(): Array<ColorWhen>;
  setColorWhenList(value: Array<ColorWhen>): void;
  addColorWhen(value?: ColorWhen, index?: number): ColorWhen;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Key.AsObject;
  static toObject(includeInstance: boolean, msg: Key): Key.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Key, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Key;
  static deserializeBinaryFromReader(message: Key, reader: jspb.BinaryReader): Key;
}

export namespace Key {
  export type AsObject = {
    name: string,
    type: string,
    layout: string,
    color?: Color.AsObject,
    maxWidth: number,
    colorWhenList: Array<ColorWhen.AsObject>,
  }
}

export class ColorWhen extends jspb.Message {
  getMatchValue(): string;
  setMatchValue(value: string): void;

  hasColor(): boolean;
  clearColor(): void;
  getColor(): Color | undefined;
  setColor(value?: Color): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ColorWhen.AsObject;
  static toObject(includeInstance: boolean, msg: ColorWhen): ColorWhen.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: ColorWhen, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ColorWhen;
  static deserializeBinaryFromReader(message: ColorWhen, reader: jspb.BinaryReader): ColorWhen;
}

export namespace ColorWhen {
  export type AsObject = {
    matchValue: string,
    color?: Color.AsObject,
  }
}

