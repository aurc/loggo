import {Component, Inject, OnInit} from '@angular/core';
import {DOCUMENT, LocationStrategy} from "@angular/common";
import {WebsocketService} from "../../services/websocket.service";
import ArrayStore from "devextreme/data/array_store";
import DataSource from "devextreme/data/data_source";

@Component({
  selector: 'app-log-stream',
  templateUrl: './log-stream.component.html',
  styleUrls: ['./log-stream.component.scss']
})
export class LogStreamComponent implements OnInit {
  private messages: any[] = [];
  dataSource: DataSource;
  private readonly messageStore: ArrayStore

  constructor(@Inject(DOCUMENT) private readonly document: Document,
              private readonly ws: WebsocketService,
              private readonly locationStrategy: LocationStrategy) {
    this.messageStore = new ArrayStore({
      key: "position",
      data: this.messages,
      // Other ArrayStore properties go here
    })
    this.dataSource = new DataSource({
      store: this.messageStore,
      reshapeOnPush: true,
      // Other DataSource properties go here
    });
  }

  ngOnInit(): void {
    this.startStream()
  }

  startStream() {
    const address = 'ws://' + this.document.location.host + '/stream'
    console.log()
    this.ws.connect(address)
    this.ws.incoming.subscribe(next => {
      const msg = {
        position: next.position,
        payload: JSON.parse(next.payload)
      }
      this.messageStore.push([{type: 'insert', data: msg}])
      console.log(msg.payload)
    })
    this.ws.send({
      "startFrom": 0
    })
  }

  stopStream() {

  }

}
