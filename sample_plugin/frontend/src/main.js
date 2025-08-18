import { defineCustomElement } from 'vue';
import ChartPluginComponent from './PluginComponent.vue';

// 1. Wrap the Vue component into a Custom Element constructor.
const ChartPluginElement = defineCustomElement(ChartPluginComponent);

// 2. Register the custom element with the browser.
// The tag name MUST contain a hyphen.
customElements.define('super-chart-plugin', ChartPluginElement);