import { defineCustomElement } from "vue";
import JSONLinterComponent from "./JSONLinterComponent.vue";

// 1. Wrap the Vue component into a Custom Element constructor.
const JSONLinterElement = defineCustomElement(JSONLinterComponent);

// 2. Register the custom element with the browser.
// The tag name MUST contain a hyphen and match the CustomElementTag in the Go plugin.
customElements.define("json-linter-formatter", JSONLinterElement);
