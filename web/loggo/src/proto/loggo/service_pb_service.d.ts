// package: loggo
// file: loggo/service.proto

import * as loggo_service_pb from "../loggo/service_pb";
import * as loggo_types_pb from "../loggo/types_pb";
import {grpc} from "@improbable-eng/grpc-web";

type LoggoLogStream = {
  readonly methodName: string;
  readonly service: typeof Loggo;
  readonly requestStream: false;
  readonly responseStream: true;
  readonly requestType: typeof loggo_types_pb.LogEntryRequest;
  readonly responseType: typeof loggo_types_pb.LogEntryResponse;
};

export class Loggo {
  static readonly serviceName: string;
  static readonly LogStream: LoggoLogStream;
}

export type ServiceError = { message: string, code: number; metadata: grpc.Metadata }
export type Status = { details: string, code: number; metadata: grpc.Metadata }

interface UnaryResponse {
  cancel(): void;
}
interface ResponseStream<T> {
  cancel(): void;
  on(type: 'data', handler: (message: T) => void): ResponseStream<T>;
  on(type: 'end', handler: (status?: Status) => void): ResponseStream<T>;
  on(type: 'status', handler: (status: Status) => void): ResponseStream<T>;
}
interface RequestStream<T> {
  write(message: T): RequestStream<T>;
  end(): void;
  cancel(): void;
  on(type: 'end', handler: (status?: Status) => void): RequestStream<T>;
  on(type: 'status', handler: (status: Status) => void): RequestStream<T>;
}
interface BidirectionalStream<ReqT, ResT> {
  write(message: ReqT): BidirectionalStream<ReqT, ResT>;
  end(): void;
  cancel(): void;
  on(type: 'data', handler: (message: ResT) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'end', handler: (status?: Status) => void): BidirectionalStream<ReqT, ResT>;
  on(type: 'status', handler: (status: Status) => void): BidirectionalStream<ReqT, ResT>;
}

export class LoggoClient {
  readonly serviceHost: string;

  constructor(serviceHost: string, options?: grpc.RpcOptions);
  logStream(requestMessage: loggo_types_pb.LogEntryRequest, metadata?: grpc.Metadata): ResponseStream<loggo_types_pb.LogEntryResponse>;
}

