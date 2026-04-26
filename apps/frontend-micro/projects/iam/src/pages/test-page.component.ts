import { Component } from '@angular/core';

@Component({
  selector: 'app-test-page',
  imports: [],
  template: `
    <div class="p-8 bg-slate-900 text-white rounded-2xl shadow-2xl border border-slate-700 m-4 animate-in fade-in zoom-in duration-500">
      <h1 class="text-3xl font-bold bg-linear-to-r from-blue-400 to-purple-500 bg-clip-text text-transparent">
        IAM Test Page 123
      </h1>
      <p class="mt-4 text-slate-400 leading-relaxed text-lg">
        Chào mừng bạn đến với trang thử nghiệm của Microfrontend IAM. 
        Trang này được render từ project <span class="text-blue-400 font-mono">iam</span> 
        và được host bởi project <span class="text-purple-400 font-mono">shell</span>.
      </p>
      
      <div class="mt-8 flex gap-4">
        <div class="p-4 bg-slate-800 rounded-xl border border-slate-700 flex-1">
          <span class="block text-sm font-semibold text-slate-500 uppercase tracking-wider">Project</span>
          <span class="text-xl font-medium">IAM (Remote)</span>
        </div>
        <div class="p-4 bg-slate-800 rounded-xl border border-slate-700 flex-1">
          <span class="block text-sm font-semibold text-slate-500 uppercase tracking-wider">Status</span>
          <span class="text-xl font-medium text-green-400 flex items-center gap-2">
            <span class="relative flex h-3 w-3">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
            </span>
            Connected
          </span>
        </div>
      </div>
    </div>
  `,
  styles: `
    :host {
      display: block;
    }
  `
})
export default class TestPageComponent { }
