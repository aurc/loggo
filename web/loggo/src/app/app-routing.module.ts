import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {DxButtonModule, DxDataGridModule, DxFormModule} from 'devextreme-angular';
import {LogStreamComponent} from './pages/log-stream/log-stream.component';
import {TemplateComponent} from './pages/template/template.component';

const routes: Routes = [
  {
    path: 'pages/template',
    component: TemplateComponent,
  },
  {
    path: 'pages/log-stream',
    component: LogStreamComponent,
  },
  {
    path: '**',
    redirectTo: 'pages/log-stream'
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes, {useHash: true}), DxDataGridModule, DxFormModule, DxButtonModule],
  exports: [RouterModule],
  declarations: [
    LogStreamComponent,
    TemplateComponent
  ]
})
export class AppRoutingModule {
}
