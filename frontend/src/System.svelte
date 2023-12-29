<script>
  import {
    CreateSigs,
    GetActiveChain,
    AddSystem,
    ConnectSystem,
  } from "../wailsjs/go/main/App";
  import { main } from "../wailsjs/go/models";

  export let systemName;
  export let sigs;
  export let comments;
  export let connectsTo;
  export let connectedTo;
  export let currentSystem;

  let viewingComments = false;
  let viewingSigs = false;
  let creatingBranch = false;

  function revertViewingState() {
    viewingComments = false;
    viewingSigs = false;
    creatingBranch = false;
  }

  function viewComments() {
    viewingComments = true;
  }

  function viewSigs() {
    viewingSigs = true;
  }

  function createBranch() {
    viewingComments = false;
    viewingSigs = false;
    creatingBranch = true;
  }

  async function createSig() {
    const text = document.getElementById("sigs");

    await CreateSigs(text.value);
    const chain = await GetActiveChain();

    sigs = chain.systems.find((system) => system.name === systemName).sigs;

    revertViewingState();
  }

  function createComment() {}

  async function createSystem() {
    try {
      creatingBranch = true;
      viewingComments = false;
      viewingSigs = false;

      const system = new main.System();

      const sn = document.getElementById("systemName");

      system.name = sn.value;

      await AddSystem(system);

      let chain = await GetActiveChain();
      const oldSystem = chain.systems.find(
        (system) => system.name === systemName
      );

      revertViewingState();

      await ConnectSystem(system.name, oldSystem.name);

      chain = await GetActiveChain();

      connectsTo = chain.systems.find(
        (system) => system.name === systemName
      ).connects_to;

      connectedTo = chain.systems.find(
        (system) => system.name === systemName
      ).connected_to;
    } catch (err) {
      console.error(err);

      alert("There was an error creating the chain! Please try again.");
    }
  }

  async function switchSystem(systemName) {
    try {
      const chain = await GetActiveChain();

      console.log(chain);

      const system = chain.systems.find((system) => system.name === systemName);

      currentSystem = system;
    } catch (err) {
      console.error(err);

      alert("There was an error switching systems! Please try again.");
    }
  }

  function drawLine(element1Id, element2Id) {
    new LeaderLine(
      document.getElementById(element1Id),
      document.getElementById(element2Id),
      {
        color: "white",
      }
    );
  }
</script>

<main>
  {#if !viewingComments && !viewingSigs && !creatingBranch}
    {#if connectsTo && !viewingComments && !viewingSigs && !creatingBranch}
      <div class="flex">
        <div class="m-auto">
          {#each connectsTo as systemName, i}
            <button
              on:click={switchSystem(systemName)}
              id={i}
              class="btn btn-primary h-5 w-5 mx-2 my-3">{systemName}</button
            >
          {/each}
        </div>
      </div>
    {/if}
    <div class="card w-96 bg-neutral text-neutral-content">
      <div class="card-body items-center text-center">
        <h2 id="systemCard" class="card-title">{systemName}</h2>
        <div class="card-actions justify-end">
          <button on:click={createBranch} class="btn btn-primary">Branch</button
          >
          <button on:click={viewComments} class="btn btn-secondary"
            >Comments</button
          >
          <button on:click={viewSigs} class="btn btn-secondary">Sigs</button>
        </div>
      </div>
    </div>
    {#if connectedTo && !viewingComments && !viewingSigs && !creatingBranch}
      <div class="flex">
        <div class="m-auto">
          <button
            on:click={switchSystem(connectedTo)}
            id={connectedTo}
            class="btn btn-primary h-5 w-5 mx-2 my-3">{connectedTo}</button
          >
        </div>
      </div>
    {/if}
  {/if}

  {#if creatingBranch}
    <div class="flex">
      <div class="m-auto">
        <input
          id="systemName"
          type="text"
          class="input input-bordered input-primary my-5"
          placeholder="System Name"
        />

        <button
          on:click={createSystem}
          class="btn btn-primary"
          id="submitSystem">Create System</button
        >
      </div>
    </div>
  {/if}

  {#if viewingSigs && !viewingComments}
    {#if sigs !== null}
      <div class="flex">
        <div class="m-auto">
          <h1 class="text-center font-bold text-xl mb-2">
            Signatures in {systemName}
          </h1>

          <ul class="mb-2">
            {#each sigs as { id, name, type }, i}
              <li>
                {id} - {name}: {type}
              </li>
            {/each}
          </ul>
        </div>
      </div>
    {:else}
      <div class="flex">
        <div class="m-auto my-2">
          <p>There are no sigs at this time.</p>
        </div>
      </div>
    {/if}
    <div class="flex">
      <div class="m-auto mb-5">
        <button on:click={revertViewingState} class="btn btn-primary"
          >Back</button
        >
      </div>
    </div>

    <div class="flex">
      <div class="m-auto">
        <h1 class="text-center font-bold mb-1">Create Sigs</h1>

        <textarea
          class="textarea textarea-primary"
          id="sigs"
          placeholder="Paste Sigs"
        ></textarea>
      </div>
    </div>

    <div class="flex">
      <div class="m-auto mb-3">
        <button on:click={createSig} class="btn btn-primary" id="submitSystem"
          >Create New Sig</button
        >
      </div>
    </div>
  {/if}

  {#if !viewingSigs && viewingComments}
    {#if comments !== null}
      <ul>
        {#each comments as { comment }, i}
          <li>
            {comment}
          </li>
        {/each}
      </ul>
    {:else}
      <p>There are no comments at this time.</p>
    {/if}

    <button on:click={revertViewingState} class="btn btn-primary">Back</button>
  {/if}
</main>
