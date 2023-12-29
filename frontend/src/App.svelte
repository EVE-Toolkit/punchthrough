<script>
  import { loop_guard } from "svelte/internal";
  import {
    ImportChain,
    SaveChain,
    NewChain,
    GetActiveChain,
  } from "../wailsjs/go/main/App";
  import { main } from "../wailsjs/go/models";
  import System from "./System.svelte";

  let chainStarted = false;
  let creatingChain = false;
  let currentSystem = null;

  async function getChain() {
    return new Promise(async (resolve, reject) => {
      try {
        const chain = await GetActiveChain();

        resolve(chain);
      } catch (err) {
        console.error(err);

        alert("There was an error fetching the chain. Please try again");

        reject(err);
      }
    });
  }

  function startNewChain() {
    creatingChain = true;
  }

  async function createChain() {
    try {
      const system = new main.System();

      const systemName = document.getElementById("systemName");

      system.name = systemName.value;

      await NewChain(system);

      creatingChain = false;
      chainStarted = true;
    } catch (err) {
      console.error(err);

      alert("There was an error creating the chain! Please try again.");
    }
  }

  async function importChain() {
    try {
      await ImportChain();

      chainStarted = true;
    } catch (err) {
      console.error(err);

      alert("There was an error importing the chain! Please try again.");
    }
  }
</script>

<main>
  <div class="flex h-screen">
    <div class="m-auto">
      {#if !chainStarted && !creatingChain}
        <div>
          <button
            on:click={startNewChain}
            class="btn btn-primary"
            id="new-chain">New Chain</button
          >
          <button
            on:click={importChain}
            class="btn btn-secondary"
            id="import-chain">Import Chain</button
          >
        </div>
      {/if}

      {#if !chainStarted && creatingChain}
        <input
          id="systemName"
          type="text"
          class="input input-bordered input-primary"
          placeholder="Root System Name"
        />

        <button on:click={createChain} class="btn btn-primary" id="submitSystem"
          >Create Chain</button
        >
      {/if}

      {#if chainStarted && currentSystem === null}
        {#await getChain()}
          <p class="font-bold text-xl">Loading chain...</p>
        {:then chain}
          {#each chain.systems as system, i}
            {#if i === 0}
              <System
                systemName={system.name}
                sigs={system.sigs}
                comments={system.comments}
                connectsTo={system.connects_to}
                connectedTo={system.connected_to}
                bind:currentSystem
              ></System>
            {/if}
          {/each}
        {:catch error}
          <p>Error loading chain: {error}</p>
        {/await}
      {/if}

      {#if chainStarted && currentSystem !== null}
        <System
          systemName={currentSystem.name}
          sigs={currentSystem.sigs}
          comments={currentSystem.comments}
          connectsTo={currentSystem.connects_to}
          connectedTo={currentSystem.connected_to}
          bind:currentSystem
        />
      {/if}
    </div>
  </div>
</main>
