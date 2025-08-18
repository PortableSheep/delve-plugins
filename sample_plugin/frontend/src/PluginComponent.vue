<template>
    <div class="chart-container">
        <h3>Super Chart Plugin</h3>
        <p>This UI is a self-contained Web Component with persistent state!</p>
        <div class="stats">
            <div class="stat">
                <strong>Chart Value:</strong>
                <div class="bar" :style="{ width: barWidth + '%' }">
                    {{ barWidth }}%
                </div>
            </div>
            <div class="stat"><strong>Updates:</strong> {{ updateCount }}</div>
            <div class="stat">
                <strong>UI Switches:</strong> {{ switchCount }}
            </div>
        </div>
        <button @click="updateChart">Update Data</button>
        <div class="status">{{ status }}</div>
    </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, getCurrentInstance } from "vue";
import delveSDK from "./delve-sdk.js";

const barWidth = ref(30);
const status = ref("Plugin loaded - persistent state");
const switchCount = ref(0);
const updateCount = ref(0);

const updateChart = () => {
    barWidth.value = Math.floor(Math.random() * 81) + 20; // Random width 20-100
    updateCount.value++;
    status.value = `Chart updated ${updateCount.value} times (switched ${switchCount.value} times)`;
};

// Example of persistent background process
const startAutoUpdate = () => {
    delveSDK.setManagedInterval(() => {
        updateChart();
        status.value = `Auto-updated at ${new Date().toLocaleTimeString()} (switched ${switchCount.value} times)`;
    }, 5000);
};

// Custom cleanup logic - only called on app shutdown
delveSDK.onCleanup(() => {
    console.log("Sample plugin shutting down - app is closing");
    status.value = "Plugin shutting down - app closing...";
    // Save any important state before shutdown
    console.log(
        `Final stats: ${updateCount.value} updates, ${switchCount.value} switches`,
    );
});

onMounted(() => {
    switchCount.value++;
    console.log(`Sample plugin UI mounted (switch #${switchCount.value})`);
    status.value = `Plugin UI loaded (switch #${switchCount.value}) - state preserved!`;

    // Setup SDK integration with Vue
    const instance = getCurrentInstance();
    delveSDK.setupVueIntegration(instance);

    // Only start auto-update on first mount to avoid multiple intervals
    if (switchCount.value === 1) {
        startAutoUpdate();
        console.log("Started background auto-update process");
    }
});

onBeforeUnmount(() => {
    console.log(
        "Sample plugin UI unmounted - but plugin continues running in background",
    );
    // Note: We don't call delveSDK.cleanup() here because the plugin should persist
    // Cleanup only happens when the entire app shuts down
});
</script>

<style>
/*
  NOTE: Styles are NOT scoped.
  They are injected into the Shadow DOM, so they won't leak out,
  but they are not scoped with data attributes like in normal Vue components.
*/
.chart-container {
    border: 1px solid #007acc;
    padding: 1rem;
    border-radius: 8px;
    background-color: #f0f5fa;
    font-family: sans-serif;
}
.stats {
    margin: 1rem 0;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
}
.stat {
    display: flex;
    align-items: center;
    gap: 0.5rem;
}
.stat strong {
    min-width: 80px;
    font-size: 0.9rem;
}
.bar {
    background-color: #007acc;
    color: white;
    padding: 5px;
    border-radius: 4px;
    transition: width 0.5s ease-in-out;
    text-align: center;
    min-width: 50px;
    flex: 1;
}
.status {
    font-size: 0.8rem;
    color: #666;
    margin-top: 1rem;
    padding: 0.5rem;
    background-color: #e8f4f8;
    border-radius: 4px;
}
</style>
