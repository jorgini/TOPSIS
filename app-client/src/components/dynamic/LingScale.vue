<script>
  import NumberForm from "./NumberForm.vue";
  import IntervalForm from "./IntervalForm.vue";
  import T1FSForm from "./T1FSForm.vue";
  import AIFSForm from "./AIFSForm.vue";
  import IT2FSForm from "./IT2FSForm.vue";
  export default {
    components: {T1FSForm, IntervalForm, NumberForm},
    data() {
      return {
        isDefaultLing: true,
        t1fs: {
          vert: [null, null, null]
        }
      }
    },
    methods: {
      chooseLingScale(event) {

      },
      correctInput(val) {
        console.log(val);
      },
      incorrectInput() {
        console.log("incorrect input");
      }
    },
  }
</script>

<template>
  <div class="ling-scale row-cols-2">
    <div class="col-4">
      <p>Лингвистическая шкала:</p>
    </div>
    <div class="col-8">
      <div class="select">
        <button type="button" :class="{opt:true, non_active: !isDefaultLing, is_active: isDefaultLing}"
                @click="chooseLingScale">По умолчанию</button>
        <button type="button" :class="{opt: true, non_active: isDefaultLing, is_active: !isDefaultLing}"
                @click="chooseLingScale">Своя конфигурация</button>
      </div>

      <p>Тип оценкки:</p>
      <select>
        <option>Число</option>
        <option>Интервал</option>
        <option>T1FS</option>
        <option>AIFS</option>
        <option>IT2FS</option>
      </select>
    </div>
  </div>

  <div class="row-cols-2">
    <div class="col-4"></div>
    <div class="col-8">
      <table class="table table-bordered">
        <tbody>
        <tr>
          <td class="mark">Плохо</td>
          <td class="eval">
            <T1FSForm v-model="t1fs" @enter-eval="correctInput" @enter-incorrect="incorrectInput"></T1FSForm>
          </td>
        </tr>
        <tr>
          <td class="mark">Плохо</td>
          <td class="eval">
            <T1FSForm v-model="t1fs" @enter-eval="correctInput" @enter-incorrect="incorrectInput"></T1FSForm>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
  @import "../../style.css";
  @import "../../assets/tasksettings.css";

  .col-8 > p {
    margin-left: 3vmin;
  }

  select {
    background-color: black;
    font-family: "Inria Sans", sans-serif;
    font-size: 2.5vmin;
    font-weight: 700;
    color: #ABF8F4;
    border-radius: 1em;
    border: 0;
    width: fit-content;
    margin-top: 1.5vmin;
  }

  table {
    border-collapse: separate;
    border-spacing: 0;
    border: 2px solid black;
    border-radius: 10px;
    width: max-content;
    overflow: hidden;
  }

  table tr {
    display: flex;
    flex-direction: row;
  }

  table td {
    background-color: #B8FCF8;
    display: flex;
    justify-content: center;
  }

  table .mark {
    min-width: 30%;
  }
</style>