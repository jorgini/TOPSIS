<script setup>
  const alt = defineModel('alt');
  const role = defineModel('role');
  const emits = defineEmits(["delete-alt"]);
</script>

<script>
  export default {
    emits: ['delete-alt'],
    data() {
      return {
        isValidTitle: true,
        isVisitableDisc: false
      }
    },
    methods: {
      validate() {
        console.log(this.role);
        this.isValidTitle = this.alt.title.length > 0 && this.alt.title.length < 101;
      },
      deleteAlt() {
        if (this.role === 'expert')
          return
        this.$emit('delete-alt', null);
      },
      switchDisc() {
        this.isVisitableDisc = !this.isVisitableDisc
      }
    }
  }
</script>

<template>
  <div class="alt">
    <div class="row-cols-3">
      <div class="col-3">
        <p>Название:</p>
      </div>
      <div class="col-6">
        <input type="text" :class="{field: true, invalid: !isValidTitle}" name="title" :readonly="role==='expert'"
               placeholder="title" maxlength="100" v-model="alt.title" @input="validate" required/>
      </div>
      <div class="col-3 right-col">
        <img alt="" src="/cancel.png" class="cancel" @click="deleteAlt">
      </div>
    </div>
    <div class="row-cols-3">
      <div class="col-3">
        <img alt="" src="/arrow.png" :class="{flag: true, close: !isVisitableDisc}" @click="switchDisc">
        <p>Описание:</p>
      </div>
      <div class="col-6">
        <textarea type="text" :class="{field: true, invisible: !isVisitableDisc}" name="description"
                  :readonly="role==='expert'" placeholder="description" maxlength="1000" v-model="alt.description"/>
      </div>
      <div class="col-3"></div>
    </div>
  </div>
</template>

<style scoped>
  @import "../../style.css";
  @import "../../assets/alts.css";
</style>